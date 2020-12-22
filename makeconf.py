#!/usr/bin/env python3
import yaml
import subprocess
import os
import fnmatch
from sys import exit
import string
import random


def generate_key(path_keyfile):
    pw = ""
    args = ["ssh-keygen", "-f", path_keyfile, "-t", "ecdsa", "-b", "521", "-N", pw]
    return subprocess.run(args)



def edit_keyname(path_keyfile):
    while os.path.isfile(path_keyfile):
        letters = string.ascii_letters
        rnd = "".join(random.choices(letters, k=4))
        path_keyfile = str(path_keyfile + rnd)
    return path_keyfile


def check_keys(expected, path_keydir):
    found = len(fnmatch.filter(os.listdir(path_keydir), '*.pub'))
    print("Found " + str(found) + "/" + str(expected) + " keys")
    kgc_msg = """
    - Not enough keys -
    Select an option for participant SSH keys (default: 1)
      1. Generate missing automatically
      2. Correct manually
    """
    if found < expected:
        kc = 0
        while kc != 1 and kc != 2:
            kc = int(input(str(kgc_msg)) or 1)
            if kc == 1:
                for x in range(found + 1, expected + 1):
                    path_keyfile = path_keydir + "/student_key" + str(x)
                    # If a key with that name exists, tack some random letters on to it
                    path_keyfile = edit_keyname(path_keyfile)
                    generate_key(path_keyfile)
                return True
            elif kc == 2:
                return False
            else:
                print("Please select a valid option")
    elif found > expected:
        bc = str(input("The keys will be used in alphabetical order (break type: \"b\")"))
        if bc == "b":
            return False
        else:
            return True
    else:
        return True


def get_byoip_info()
    args = ["aws", "ec2", "describe-public-ipv4-pools", "--query", "PublicIpv4Pools[].PoolId"]
    ipv4id = subprocess.run(args)

if __name__ == "__main__":
    default_domain = "examples.nu"
    domainlist = []
    student_key = []
    student_keyname = []
    keydir = "keys"
    parameterdir = "parameters"

    project = str(input("Name of the project: (default: DNScourse)") or "DNScourse")
    prefix = str(input("What name prefix? (default: group) ") or "group")

    if os.path.isdir(keydir):
        pass
    else:
        try:
            os.mkdir(keydir)
        except Exception as e:
            print(e)

    if os.path.isdir(parameterdir):
        pass
    else:
        try:
            os.mkdir(parameterdir)
        except Exception as e:
            print(e)

    nr = int(input("How many lab environments do you want? (default: 1) ") or 1)
    kc_msg = """
    Select an option for participant SSH keys (default: 1)
      1. shared key (generate)
      2. shared key (pre-generated)
      3. unique keys (generate)
      4. unique keys (pre-generated)

    """
    keychoice = int(input(str(kc_msg)) or 1)
    if keychoice == 1:
        print("*generates shared key*")
        keyfile = keydir + "/student_key"
        keyfile = edit_keyname(keyfile)
        generate_key(keyfile)
        for i in range(nr):
            with open(keyfile + ".pub", "r") as file:
                student_key.append(file.read())
                student_keyname.append(str(keyfile))
    elif keychoice == 2:
        print("*asks for shared key*")
        while True:
            if check_keys(1, keydir):
                try:
                    for file in os.listdir(keydir):
                        if file.endswith(".pub"):
                            for i in range(nr):
                                key = open(os.path.join(keydir, file)).read()
                                student_key.append(key)
                                student_keyname.append(str(keydir+"/"+file))
                            break
                except Exception as e:
                    print(e)
                break
            else:
                choice = str(input("Press enter to check your keys again or type \"e\" to exit script "))
                if choice == "e":
                    exit()
    elif keychoice == 3:
        print("*generates " + str(nr) + " keys*")
        for i in range(1, nr + 1):
            keyfile = keydir + "/student_key" + str(i)
            keyfile = edit_keyname(keyfile)
            generate_key(keyfile)
            with open(keyfile + ".pub", "r") as file:
                student_key.append(file.read())
                student_keyname.append(str(keyfile))
    elif keychoice == 4:
        print("*asks for " + str(nr) + " keys*")
        while True:
            if check_keys(nr, keydir):
                try:
                    for file in os.listdir(keydir):
                        if file.endswith(".pub"):
                            key = open(os.path.join(keydir, file)).read()
                            student_key.append(key)
                            student_keyname.append(str(keydir+"/"+file))
                except Exception as e:
                    print(e)
                break
            else:
                choice = str(input("Press enter to check your keys again or type \"e\" to exit script "))
                if choice == "e":
                    exit()
    else:
        print("*generates shared key*")
        keyfile = keydir + "/student_key"
        keyfile = edit_keyname(keyfile)
        generate_key(keyfile)
        for i in range(1, nr + 1):
            with open(keyfile + ".pub", "r") as file:
                student_key.append(file.read())
                student_keyname.append(str(keyfile))

    admin_key = str(input("Name of your admin key (default: DNSkurs) ") or "DNSkurs")

    domchoice = str(input("Use a common domain name? (default: yes) ") or "yes")
    if domchoice == "yes":
        domain = str(input("Which domain? (default: " + default_domain + ") ") or default_domain)
        for i in range(1, nr + 1):
            domainlist.append(prefix+str(i)+"."+domain)
    else:
        for i in range(1, nr + 1):
            domain = str(input("Enter domain " + str(i) + "(of " + str(nr) + ") ") or default_domain)
            domainlist.append(domain)

    Ipv4PoolLAN = 0
    Ipv4PoolSNS = 0
    while Ipv4PoolLAN != 1 and Ipv4PoolLAN != 2 and Ipv4PoolSNS != 1 and Ipv4PoolSNS != 2:
        Ipv4PoolLAN_msg = """
        Select an option for the Ipv4 pool used in the Virtual LAN (default: 1)
        1. Amazon provided Ipv4 pool
        2. Bring your own Ipv4 pool
        """
        Ipv4PoolSNS_msg = """
        Select an option for the Ipv4 pool used in the Secondary nameserver(default: 1)
        1. Amazon provided Ipv4 pool
        2. Bring your own Ipv4 pool
        """
        Ipv4PoolLAN = int(input(str(Ipv4PoolLAN_msg)) or 1)
        Ipv4PoolSNS = int(input(str(Ipv4PoolSNS_msg)) or 1)
        if Ipv4PoolLAN == 1 and Ipv4PoolSNS == 1:
            BYOIP = ""
        elif Ipv4PoolLAN == 2 or Ipv4PoolSNS == 2:
            BYOIP = str(input("Id of your Ipv4 address pool: "))
        else:
            print("Please select valid options")

    Ipv6PoolLAN = 0
    Ipv6PoolSNS = 0
    while Ipv6PoolLAN != 1 and Ipv6PoolLAN != 2 and Ipv6PoolSNS != 1 and Ipv6PoolSNS != 2:
        Ipv6PoolLAN_msg = """
        Select an option for the Ipv6 pool used in the Virtual LAN (default: 1)
        1. Amazon provided Ipv6 pool
        2. Bring your own Ipv6 pool
        """
        Ipv6PoolSNS_msg = """
        Select an option for the Ipv6 pool used in the Secondary nameserver (default: 1)
        1. Amazon provided Ipv6 pool
        2. Bring your own Ipv6 pool
        """
        Ipv6PoolLAN = int(input(str(Ipv6PoolLAN_msg)) or 1)
        Ipv6PoolSNS = int(input(str(Ipv6PoolSNS_msg)) or 1)
        if Ipv6PoolLAN == 1 and Ipv6PoolSNS == 1:
            BYOIPv6 = ""
            Ipv6CidrBlock = ""
        elif Ipv6PoolLAN == 2 or Ipv6PoolSNS == 2:
            BYOIPv6 = str(input("Id of your Ipv6 address pool: "))
            Ipv6CidrBlock = str(input("/56 cidr block of your own Ipv6 address pool: "))
        else:
            print("Please select valid options")


    paramfile = parameterdir+"/"+str(input("Name of parameter file: "))+".yaml"
    while os.path.isfile(paramfile):
        fc_msg = """
        - Filename already exists -
        Select an option for filename (default: 1)
        1. Choose a different filename
        2. Write over existing file
        """
        fc = int(input(fc_msg) or 1)
        if fc == 1:
            paramfile = parameterdir+"/"+str(input("Name of parameter file: "))+".yaml"
        elif fc == 2:
            break
        else:
            paramfile = parameterdir+"/"+str(input("Name of parameter file: "))+".yaml"

    paramlist = []
    for i in range(nr):
        paramlist.append({"Project": project, "Prefix": prefix, "GroupName": prefix + str(i + 1),
                          "DomainName": domainlist[i], "KeyName": admin_key,
                          "StudentKey": student_key[i], "StudentKeyName": student_keyname[i], "Ipv4PoolLAN": Ipv4PoolLAN, "Ipv4PoolSNS": Ipv4PoolSNS,
                          "BYOIP": BYOIP, "Ipv6PoolLAN": Ipv6PoolLAN, "Ipv6PoolSNS": Ipv6PoolSNS, "BYOIPv6": BYOIPv6, "Ipv6CidrBlock": Ipv6CidrBlock})
    with open(paramfile, "w") as file:
        yaml.dump(paramlist, file)
