hackernews.ir
ip: 10.0.0.1


hn1:3030
hn-summary:3040
mysql:3306

=== VIRTUAL PRIVATE SERVER
ssh root@10.0.0.1

apt install fail2ban

ufw allow port 22
ufw allow port 3030
ufw enable

1 server: 2 service with 2 port

nginx
2. http server 
1. reverse proxy
3. load balancer

hackernews.ir -> 10.0.0.1:3030
summary.hackernews.ir -> 10.0.0.1:3040


hackernews.ir