from docker import Client
import requests, json

# define program-wide variables
BIGIP_ADDRESS = '[Address of BIG-IP]'
BIGIP_USER = '[Admin User]'
BIGIP_PASS = '[Admin Password]'

DOCKER_HOSTS = ['[List of Docker Hosts]']
HTTP_PORT = '80'
HTTP_PROTOCOL = '%s/tcp' %(HTTP_PORT)

#requests.packages.urllib3.disable_warnings()

clients = []
pools = {}
data_group = {}
#
# functions
#

def create_pool(bigip, name, members):
        payload = {}

        # convert member format
        payload_members = [ { 'name' : member } for member in members ]

        # define test pool
        payload['name'] = name
        payload['description'] = 'built by docker_to_f5_bigip.py'
        payload['loadBalancingMode'] = 'least-connections-member'
        payload['monitor'] = 'http'
        payload['members'] = members
        req = bigip.post('%s/ltm/pool' % BIGIP_URL_BASE, data=json.dumps(payload))

def update_pool(bigip, name, members):
        payload = {}

        # convert member format
        payload_members = [ { 'name' : member } for member in members ]

        # define test pool
        payload['name'] = name
        payload['members'] = members
        req = bigip.patch('%s/ltm/pool/%s' % (BIGIP_URL_BASE, name) , data=json.dumps(payload))



#update DataGroup
def update_dg(bigip, name, data_group):
        payload = {}

        payload['records'] =  [{'data':r[1],'name':r[0]} for r in  data_group.items()]

        req = bigip.patch('%s/ltm/data-group/internal/%s' % (BIGIP_URL_BASE, name), data=json.dumps(payload))


#
# connect to docker hosts
#
for host in DOCKER_HOSTS:
    try:
        cli = Client(host)
        cli.info()
    except:
        print "failled to connect to",host
        continue
    clients.append(cli)

containers = {}
#
# grab info about containers
#

for cli in clients:
    tmp = [ c['Id'] for c in cli.containers()]
    for cid in tmp:
        details = cli.inspect_container(cid)
        containers[cid[:12]] = {'Name': details['Name'][1:],
                           'IPv4': details['NetworkSettings']['IPAddress'],
                           'Ports': details['NetworkSettings']['Ports'].keys(),
                       }
#
# build list of HTTP services
#
for cnt in containers.values():
    ports = cnt['Ports']
    if HTTP_PROTOCOL in ports:
        ip_port = '%s:%s' %(cnt['IPv4'],HTTP_PORT)
        con_name = cnt['Name']
        pool_name = 'docker_%s_pool' %(con_name.split('-')[0])
        pool = pools.get(pool_name,[])
        pool.append(ip_port)
        pools[pool_name] = pool
        data_group[cnt['Name']] = ip_port

# REST resource for BIG-IP that all other requests will use
bigip = requests.session()
bigip.auth = (BIGIP_USER, BIGIP_PASS)
bigip.verify = False
bigip.headers.update({'Content-Type' : 'application/json'})

# Requests requires a full URL to be sent as arg for every request, define base URL globally here
BIGIP_URL_BASE = 'https://%s/mgmt/tm' % BIGIP_ADDRESS

#
# grab all pool names
#
req =  bigip.get('%s/ltm/pool' % BIGIP_URL_BASE)
pool_json = req.json()

pool_names = [a['name'] for a in pool_json['items'] if a['name'].startswith('docker_')]
local_pools = set(pool_names)
remote_pools = set(pools.keys())

to_delete = local_pools - remote_pools
to_add = remote_pools  - local_pools
to_update = remote_pools & local_pools

for pname in to_delete:
    req =  bigip.delete('%s/ltm/pool/%s' % (BIGIP_URL_BASE, pname))

for pname in to_add:
    payload = {}
    create_pool(bigip, pname, pools[pname])

for pname in to_update:
    payload = {}
    update_pool(bigip, pname, pools[pname])
update_dg(bigip, 'dg_docker_container', data_group)
update_dg(bigip, 'dg_docker_pool', dict((a[7:-5],a) for a in pools.keys()))
