from fabric.api import *
from fabric.contrib.files import *

env.hosts = ['root@107.170.45.49']

OPENRESTY = 'ngx_openresty-1.7.2.1'

@task
def install_packages():
	sudo('apt-get update')
	sudo('apt-get install -y libreadline-dev libncurses5-dev libpcre3-dev libssl-dev perl make')

@task
def install_openresty():
	if exists(OPENRESTY) is False:
		run('wget -O - http://openresty.org/download/%s.tar.gz | tar xzf -' % (OPENRESTY))
	with cd(OPENRESTY):
		run('./configure')
		run('make')
		sudo('make install')

@task
def nginx(command):
	sudo('service nginx %(command)s' % locals())

@task
def add_init_script():
	put('init.d/nginx', '/etc/init.d/nginx')
	sudo('chmod +x /etc/init.d/nginx')
	sudo('update-rc.d -f nginx defaults')
	nginx('start')

@task
def update_config():
	with cd('/usr/local/openresty/nginx/conf'):
		sudo('mv -n nginx.conf nginx.conf.default')
		put('conf/nginx.conf', 'nginx.conf')
	nginx('reload')

@task
def update_scripts():
	sudo('mkdir -p /usr/local/openresty/nginx/lua')
	with cd('/usr/local/openresty/nginx/lua'):
		put('lua/redirect.lua', 'redirect.lua')
		put('lua/parse_redirect.lua', 'parse_redirect.lua')
	nginx('reload')

@task
def install():
	install_packages()
	install_openresty()

@task
def init():
	add_init_script()
	update_config()
	update_scripts()
