import subprocess
import networkx as nx
import socket
import sys
import os
import re
import getopt
import json
import datetime
import pprint

import scapy
from scapy.all import *

import matplotlib
# matplotlib.use('TKAgg')
import matplotlib.pyplot as plt

import pymongo
from pymongo import MongoClient

# Needs mongodb to store targets as dictionaries'
def setup_db():
  try:
    dbconn = MongoClient()
  except:
    print ('Could not connect to db, make sure you started mongo')
  try:
    db = dbconn.test
  except:
    print ('Could not get the pentest database')
  try:
    collection = db.targets
  except:
     print ('Could not get the collection of targets')
  return collection

# Retruns all targets stored in the database
def get_all_targets():
   for t in coll.find():
      pprint.pprint(t,indent=4)

class target():
  def __init__(self,ip,coll=None):
    self.target={}
    self.ip=ip[0]
    self.target['hostname']=ip[1]
    self.collection=coll

  def __str__(self):
    return json.dumps(self.target,sort_keys=True,indent=4)
	
  def make_target(self):
    self.target_ip()
    self.port_scanned()
    self.traceroute()
    self.target['Timestamp']=str(datetime.datetime.utcnow())
    self.collection.insert(self.target)

  def target_ip(self):
    self.target['ip']=self.ip

  def port_scanned(self):
    ports=nmap_scan(self.ip)
    self.target['ports']=ports

  def traceroute(self):
    hops=scapytraceroute(self.ip)
    self.target['traceroute']=hops

  def bruteforce_reversedns(base_ip):
    ip_list=[]
    for i in range(255):
      try:
        (hostname,alias,ip)=socket.gethostbyaddr(base_ip+str(i))
        ip_list.append((ip[0],hostname))
      except:
        pass
      return ip_list

  def nmap_scan(host):
    ports = []
    cmd = 'sudo nmap -Pn -sS ' + host
    print ('Scanning: ') + cmd
    p=subprocess.Popen(cmd,shell=True,stdout=subprocess.PIPE,stderr=subprocess.PIPE)
    (pout,perr)=p.communicate()
    foobar=re.compile('tcp')
    for line in pout.split('\n'):
      if foobar.search(line):
        print (line)
        ports.append(line)
    return ports

  def localtraceroute(host,num_hops):
    hops=[]
    trace='traceroute -m %d %s' % (num_hops,host)
    print (trace)
    res=subprocess.Popen(trace,shell=True,stdout=subprocess.PIPE,stderr=subprocess.PIPE)
    (pstdout,psterr)=res.communicate()
    lines=pstdout.split('\n')
    for line in lines[:num_hops]:
       hops.append(line.split(' ')[num_hops-1].rstrip(')').lstrip('('))
    return hops

  def scapytraceroute(host):
    hops=[]
    try:
      res,unans=traceroute(host)
    except:
      print ('Could not trace route with scapy !')
    return hops
    host_key=res.get_trace().keys()[0]
    for key in res.get_trace()[host_key].keys():
      hops.append(res.get_trace()[host_key][key][0])
    return hops

  def traceroute_plot(targets):
    g=nx.Graph()
    source=socket.gethostbyname(socket.gethostname())
    for t in targets:
      hops=scapytraceroute(t)
      print (hops)
      g.add_node(t)
      g.add_edge(source,hops[0])
      if len(hops) >= 1:
        for hop in hops:
          next_hop_index=hops.index(hop)+1
          if next_hop_index != len(hops):
            g.add_edge(hop,hops[next_hop_index])
          else:
            g.add_edge(hop,t)
            nx.draw(g,with_labels=False)
            plt.savefig("/var/tmp/trace.png")
            nx.write_dot(g,"/var/tmp/trace.dot")

  def traceroute_plot_from_db(targets):
    g=nx.Graph()
    source=socket.gethostbyname(socket.gethostname())
    for t in targets:
      hops=t['traceroute']
      print (hops)
      g.add_node(t['ip'])
      g.add_edge(source,hops[0])
      if len(hops) >= 1:
        for hop in hops:
          next_hop_index=hops.index(hop)+1
          if next_hop_index != len(hops):
            g.add_edge(hop,hops[next_hop_index])
          else:
            g.add_edge(hop,t['ip'])
            nx.draw(g,with_labels=False)
            plt.savefig("/var/tmp/trace.png") 
            nx.write_dot(g,"/var/tmp/trace.dot")

  def main():
    targets=[]
    try:
      fh=open(targets_file,'r')
    except:
      print ('targets.list file not present')
      sys.exit()
    for line in fh.readlines():
      targets.append(line.strip('\n'))
    traceroute_plot(targets)
	
  def readopt():
    try:
      options, remainder = getopt.getopt(sys.argv[1:],'b:s:t:f:')
    except getopt.GetoptError as err:
      print ('option error.')
      usage()
      sys.exit(2)
    global base_ip,host_to_scan,host_to_traceroute,targets_file
    base_ip = '0.0.0.0'
    host_to_scan = '1.1.1.1'
    host_to_traceroute = '1.1.1.1'
    targets_file = 'targets.list'
    for opt, arg in options:
      if opt == '-b':
        base_ip = arg
      elif opt == '-s':
        host_to_scan = arg
      elif opt == '-t':
        host_to_traceroute == arg
      elif opt == '-f':
        targets_file == arg
      else:
        usage()
        sys.exit(2)

  def usage():
    print ('This code plots the traceroute to a set of hosts')

if __name__=="__main__":
  # readopt()
  # sys.exit(main())
  targets = target()
  targets.main()
else:
  print ('The db setup will called and can be refered as trace.coll() in an interactive shell')
  coll = setup_db()
