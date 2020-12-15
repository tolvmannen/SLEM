#!/usr/bin/env python3
import yaml
import subprocess
import argparse
import time

def VPC(params):
    Project = params["Project"]
    GroupName = params["GroupName"]
    Prefix = params["Prefix"]
    StackName = params["GroupName"] + "-VPC"
    Ipv6Pool = str(params["Ipv6PoolLAN"])
    print("Deploying "+StackName)
    args = ["aws", "cloudformation", "deploy", "--template-file", "templates/VPC.yaml", "--stack-name", StackName, "--parameter-overrides", "Ipv6Pool="+Ipv6Pool, "--tags", "Project="+Project, Prefix+"="+GroupName, "Name="+StackName]
    subprocess.run(args)

def VPC_update(params, CidrBlock):
    Ipv6CidrBlock = str(CidrBlock)
    Ipv6Pool = str(params["Ipv6PoolLAN"])
    Project = params["Project"]
    GroupName = params["GroupName"]
    Prefix = params["Prefix"]
    StackName = params["GroupName"] + "-VPC"
    print("Updating "+StackName)
    args = ["aws", "cloudformation", "deploy", "--template-file", "templates/VPC_updated.yaml", "--stack-name", StackName, "--parameter-overrides", "Ipv6Pool="+Ipv6Pool, "Ipv6CidrBlock="+Ipv6CidrBlock, "--tags", "Project="+Project, Prefix+"="+GroupName, "Name="+StackName]
    subprocess.run(args)

def VPCSlaveDNS(params):
    Project = params["Project"]
    StackName = "VPCSlaveDNS"
    Ipv6Pool = str(params["Ipv6PoolSNS"])
    print("Deploying "+StackName)
    args = ["aws", "cloudformation", "deploy", "--template-file", "templates/VPCSlaveDNS.yaml", "--stack-name", StackName, "--parameter-overrides", "Ipv6Pool="+Ipv6Pool, "--tags", "Project="+Project, "Name="+StackName]
    subprocess.run(args)

def VPCSlaveDNS_update(params, CidrBlockSlave):
    Project = params["Project"]
    StackName = "VPCSlaveDNS"
    Ipv6Pool = str(params["Ipv6PoolSNS"])
    Ipv6CidrBlock = str(CidrBlockSlave)
    print("Updating "+StackName)
    args = ["aws", "cloudformation", "deploy", "--template-file", "templates/VPCSlaveDNS_updated.yaml", "--stack-name", StackName, "--parameter-overrides", "Ipv6Pool="+Ipv6Pool, "Ipv6CidrBlock="+Ipv6CidrBlock, "--tags", "Project="+Project, "Name="+StackName]
    subprocess.run(args)

def securityGroups(params):
    Project = params["Project"]
    GroupName = params["GroupName"]
    Prefix = params["Prefix"]
    StackName = params["GroupName"] + "-SG"
    VPCStackName = params["GroupName"] + "-VPC"
    print("Deploying "+StackName)
    args = ["aws", "cloudformation", "deploy", "--template-file", "templates/SecurityGroups.yaml", "--stack-name", StackName, "--parameter-overrides", "VPCStackName="+VPCStackName, "--tags", "Project="+Project, Prefix+"="+GroupName, "Name="+StackName]
    subprocess.run(args)

def proxy(params):
    Project = params["Project"]
    GroupName = params["GroupName"]
    Prefix = params["Prefix"]
    StackName = params["GroupName"] + "-proxy"
    KeyName = params["KeyName"]
    VPCStackName = params["GroupName"] + "-VPC"
    HostName = params["GroupName"] + "-proxy"
    Ipv4Pool = params["Ipv4PoolLAN"]
    BYOIP = params["BYOIP"]
    print("Deploying "+StackName)
    args = ["aws", "cloudformation", "deploy", "--template-file", "templates/Proxy.yaml", "--stack-name", StackName, "--parameter-overrides", "VPCStackName="+VPCStackName, "KeyName="+KeyName, "HostName="+HostName, "Ipv4Pool="+str(Ipv4Pool), "BYOIP="+BYOIP, "--tags", "Project="+Project, Prefix+"="+GroupName, "Name="+StackName]
    subprocess.run(args)

def resolver(params):
    Project = params["Project"]
    GroupName = params["GroupName"]
    Prefix = params["Prefix"]
    StackName = params["GroupName"] + "-resolver"
    KeyName = params["KeyName"]
    StudentKey = params["StudentKey"]
    StudentKeyName = params["StudentKeyName"]
    VPCStackName = params["GroupName"] + "-VPC"
    HostName = params["GroupName"] + "-resolver"
    print("Deploying "+StackName)
    args = ["aws", "cloudformation", "deploy", "--template-file", "templates/Resolver.yaml", "--stack-name", StackName, "--parameter-overrides", "VPCStackName="+VPCStackName, "KeyName="+KeyName, "HostName="+HostName, "StudentKey="+StudentKey, "--tags", "Project="+Project, Prefix+"="+GroupName, "Name="+StackName]
    subprocess.run(args)

def jumpgate(params):
    Project = params["Project"]
    GroupName = params["GroupName"]
    Prefix = params["Prefix"]
    StackName = params["GroupName"] + "-jumpgate"
    KeyName = params["KeyName"]
    StudentKey = params["StudentKey"]
    StudentKeyName = params["StudentKeyName"]
    VPCStackName = params["GroupName"] + "-VPC"
    HostName = params["GroupName"] + "-jumpgate"
    print("Deploying "+StackName)
    args = ["aws", "cloudformation", "deploy", "--template-file", "templates/Jumpgate.yaml", "--stack-name", StackName, "--parameter-overrides",  "VPCStackName="+VPCStackName, "KeyName="+KeyName, "HostName="+HostName, "StudentKey="+StudentKey, "--tags", "Project="+Project, Prefix+"="+GroupName, "Name="+StackName]
    subprocess.run(args)

def LANclient(params):
    Project = params["Project"]
    GroupName = params["GroupName"]
    Prefix = params["Prefix"]
    StackName = params["GroupName"] + "-LANclient"
    KeyName = params["KeyName"]
    VPCStackName = params["GroupName"] + "-VPC"
    HostName = params["GroupName"] + "-LANclient"
    print("Deploying "+StackName)
    args = ["aws", "cloudformation", "deploy", "--template-file", "templates/LANclient.yaml", "--stack-name", StackName, "--parameter-overrides", "VPCStackName="+VPCStackName, "KeyName="+KeyName, "HostName="+HostName, "--tags", "Project="+Project, Prefix+"="+GroupName, "Name="+StackName]
    subprocess.run(args)

def DNSmaster(params):
    Project = params["Project"]
    GroupName = params["GroupName"]
    Prefix = params["Prefix"]
    StackName = params["GroupName"] + "-DNSmaster"
    KeyName = params["KeyName"]
    StudentKey = params["StudentKey"]
    StudentKeyName = params["StudentKeyName"]
    VPCStackName = params["GroupName"] + "-VPC"
    HostName = params["GroupName"] + "-DNSmaster"
    print("Deploying "+StackName)
    args = ["aws", "cloudformation", "deploy", "--template-file", "templates/DNSmaster.yaml", "--stack-name", StackName, "--parameter-overrides", "VPCStackName="+VPCStackName, "KeyName="+KeyName, "HostName="+HostName, "StudentKey="+StudentKey, "--tags", "Project="+Project, Prefix+"="+GroupName, "Name="+StackName]
    subprocess.run(args)

def webbserver(params):
    Project = params["Project"]
    GroupName = params["GroupName"]
    Prefix = params["Prefix"]
    StackName = params["GroupName"] + "-webbserver"
    KeyName = params["KeyName"]
    StudentKey = params["StudentKey"]
    StudentKeyName = params["StudentKeyName"]
    VPCStackName = params["GroupName"] + "-VPC"
    HostName = params["GroupName"] + "-webbserver"
    DomainName = params["DomainName"]
    print("Deploying "+StackName)
    args = ["aws", "cloudformation", "deploy", "--template-file", "templates/Webbserver.yaml", "--stack-name", StackName, "--parameter-overrides", "VPCStackName="+VPCStackName, "KeyName="+KeyName, "HostName="+HostName, "StudentKey="+StudentKey, "DomainName="+DomainName, "--tags", "Project="+Project, Prefix+"="+GroupName, "Name="+StackName]
    subprocess.run(args)

def mailserver(params):
    Project = params["Project"]
    GroupName = params["GroupName"]
    Prefix = params["Prefix"]
    StackName = params["GroupName"] + "-mailserver"
    KeyName = params["KeyName"]
    StudentKey = params["StudentKey"]
    StudentKeyName = params["StudentKeyName"]
    VPCStackName = params["GroupName"] + "-VPC"
    HostName = params["GroupName"] + "-mailserver"
    DomainName = params["DomainName"]
    print("Deploying "+StackName)
    args = ["aws", "cloudformation", "deploy", "--template-file", "templates/Mailserver.yaml", "--stack-name", StackName, "--parameter-overrides", "VPCStackName="+VPCStackName, "KeyName="+KeyName, "HostName="+HostName, "StudentKey="+StudentKey, "DomainName="+DomainName, "--tags", "Project="+Project, Prefix+"="+GroupName, "Name="+StackName]
    subprocess.run(args)

def DNSslave(params):
    Project = params["Project"]
    GroupName = params["GroupName"]
    Prefix = params["Prefix"]
    StackName = params["GroupName"] + "-DNSslave"
    KeyName = params["KeyName"]
    StudentKey = params["StudentKey"]
    StudentKeyName = params["StudentKeyName"]
    VPCSlaveStackName = "VPCSlaveDNS"
    HostName = params["GroupName"] + "-DNSslave"
    Ipv4Pool = params["Ipv4PoolSNS"]
    BYOIP = params["BYOIP"]
    print("Deploying "+StackName)
    args = ["aws", "cloudformation", "deploy", "--template-file", "templates/DNSslave.yaml", "--stack-name", StackName, "--parameter-overrides", "VPCSlaveStackName="+VPCSlaveStackName, "KeyName="+KeyName, "HostName="+HostName, "StudentKey="+StudentKey, "Ipv4Pool="+str(Ipv4Pool), "BYOIP="+BYOIP, "--tags", "Project="+Project, Prefix+"="+GroupName, "Name="+StackName]
    subprocess.run(args)

def deployAll(paramlist):
    VPCSlaveDNS(paramlist[0])
    if paramlist[0]["Ipv6PoolSNS"] == 2:
        VPCSlave_ID = subprocess.check_output("aws ec2 describe-vpcs --filter Name=tag:Name,Values=VPCSlaveDNS --query Vpcs[].VpcId --output text", shell=True)
        VPCSlave_ID = str(VPCSlave_ID)
        VPCSlave_ID = VPCSlave_ID.split("'")[1].strip("\\r\\n")
        Ipv6CidrBlockSlave = str(paramlist[0]["Ipv6CidrBlock"])[:15]+"0000::/56"
        subprocess.run(["aws", "ec2", "associate-vpc-cidr-block", "--no-amazon-provided-ipv6-cidr-block", "--vpc-id", VPCSlave_ID, "--ipv6-cidr-block", Ipv6CidrBlockSlave, "--ipv6-pool", paramlist[0]["BYOIPv6"]])
        print("Waiting for cidr block to associate")
        CB_state = ""
        while CB_state != "associated":
            CB_state = subprocess.check_output("aws ec2 describe-vpcs --filter Name=tag:Name,Values=VPCSlaveDNS --query Vpcs[].Ipv6CidrBlockAssociationSet[].Ipv6CidrBlockState.State --output text", shell=True)
            CB_state = str(CB_state)
            CB_state = CB_state.split("'")[1].strip("\\r\\n")
            time.sleep(5)
    else:
        Ipv6CidrBlockSlave = ""
    VPCSlaveDNS_update(paramlist[0], Ipv6CidrBlockSlave)
    for index, params in enumerate(paramlist):
        VPC(params)
        if params["Ipv6PoolLAN"] == 2:
            index = str(index+1)
            if len(index) == 1:
                index = "0"+index
            VPC_ID = subprocess.check_output("aws ec2 describe-vpcs --filter Name=tag:Name,Values="+params["GroupName"]+"-VPC --query Vpcs[].VpcId --output text", shell=True)
            VPC_ID = str(VPC_ID)
            VPC_ID = VPC_ID.split("'")[1].strip("\\r\\n")
            Ipv6CidrBlock = str(params["Ipv6CidrBlock"])[:15]+str(index)+"00::/56"
            subprocess.run(["aws", "ec2", "associate-vpc-cidr-block", "--no-amazon-provided-ipv6-cidr-block", "--vpc-id", VPC_ID, "--ipv6-cidr-block", Ipv6CidrBlock, "--ipv6-pool", params["BYOIPv6"]])
            print("Waiting for cidr block to associate")
            CB_state = ""
            while CB_state != "associated":
                CB_state = subprocess.check_output("aws ec2 describe-vpcs --filter Name=tag:Name,Values="+params["GroupName"]+"-VPC --query Vpcs[].Ipv6CidrBlockAssociationSet[].Ipv6CidrBlockState.State --output text", shell=True)
                CB_state = str(CB_state)
                CB_state = CB_state.split("'")[1].strip("\\r\\n")
                time.sleep(5)
        else:
            Ipv6CidrBlock = ""
        VPC_update(params, Ipv6CidrBlock)
        securityGroups(params)
        proxy(params)
        resolver(params)
        jumpgate(params)
        LANclient(params)
        DNSmaster(params)
        webbserver(params)
        mailserver(params)
        DNSslave(params)

def deployOne(params,index):
    VPCSlaveDNS(params)
    if params["Ipv6PoolSNS"] == 2:
        VPCSlave_ID = subprocess.check_output("aws ec2 describe-vpcs --filter Name=tag:Name,Values=VPCSlaveDNS --query Vpcs[].VpcId --output text", shell=True)
        VPCSlave_ID = str(VPCSlave_ID)
        VPCSlave_ID = VPCSlave_ID.split("'")[1].strip("\\r\\n")
        Ipv6CidrBlockSlave = str(params["Ipv6CidrBlock"])[:15]+"0000::/56"
        subprocess.run(["aws", "ec2", "associate-vpc-cidr-block", "--no-amazon-provided-ipv6-cidr-block", "--vpc-id", VPCSlave_ID, "--ipv6-cidr-block", Ipv6CidrBlockSlave, "--ipv6-pool", params["BYOIPv6"]])
        print("Waiting for cidr block to associate")
        CB_state = ""
        while CB_state != "associated":
            CB_state = subprocess.check_output("aws ec2 describe-vpcs --filter Name=tag:Name,Values=VPCSlaveDNS --query Vpcs[].Ipv6CidrBlockAssociationSet[].Ipv6CidrBlockState.State --output text", shell=True)
            CB_state = str(CB_state)
            CB_state = CB_state.split("'")[1].strip("\\r\\n")
            time.sleep(5)
    else:
        Ipv6CidrBlockSlave = ""
    VPCSlaveDNS_update(params, Ipv6CidrBlockSlave)
    VPC(params)
    if params["Ipv6PoolLAN"] == 2:
        index = str(index)
        if len(index) == 1:
            index = "0"+index
        VPC_ID = subprocess.check_output("aws ec2 describe-vpcs --filter Name=tag:Name,Values="+params["GroupName"]+"-VPC --query Vpcs[].VpcId --output text", shell=True)
        VPC_ID = str(VPC_ID)
        VPC_ID = VPC_ID.split("'")[1].strip("\\r\\n")
        Ipv6CidrBlock = str(params["Ipv6CidrBlock"])[:15]+str(index)+"00::/56"
        subprocess.run(["aws", "ec2", "associate-vpc-cidr-block", "--no-amazon-provided-ipv6-cidr-block", "--vpc-id", VPC_ID, "--ipv6-cidr-block", Ipv6CidrBlock, "--ipv6-pool", params["BYOIPv6"]])
        print("Waiting for cidr block to associate")
        CB_state = ""
        while CB_state != "associated":
            CB_state = subprocess.check_output("aws ec2 describe-vpcs --filter Name=tag:Name,Values="+params["GroupName"]+"-VPC --query Vpcs[].Ipv6CidrBlockAssociationSet[].Ipv6CidrBlockState.State --output text", shell=True)
            CB_state = str(CB_state)
            CB_state = CB_state.split("'")[1].strip("\\r\\n")
            time.sleep(5)
    else:
        Ipv6CidrBlock = ""
    VPC_update(params, Ipv6CidrBlock)
    securityGroups(params)
    proxy(params)
    resolver(params)
    jumpgate(params)
    LANclient(params)
    DNSmaster(params)
    webbserver(params)
    mailserver(params)
    DNSslave(params)

if __name__ == "__main__":
    parser = argparse.ArgumentParser(prog='deploy', description='Deploy an AWS environment from YAML-file')
    parser.add_argument('-f',
                        '--file',
                        action='store',
                        required=True,
                        help='path to YAML configuration file')
    args = parser.parse_args()
    try:
        conf_file = open(args.file)
        paramlist = yaml.load(conf_file, Loader=yaml.FullLoader)
    except:
        print('Could not open file ', str(args.file) )


    choice_msg = """
    Select an option for deployment (default: 1)
    1. Deploy all environments
    2. Deploy a specific groups environment
    3. Deploy a specific server for a specific group
    """
    choice = int(input(str(choice_msg)) or 1)
    if choice == 1:
        deployAll(paramlist)
    elif choice == 2:
        group_nr = int(input("Which group?"))
        deployOne(paramlist[group_nr-1],group_nr)
    elif choice == 3:
        group_nr = int(input("Which group?"))
        sc_msg = """
        Select an option for which server to deploy
        1. Proxy
        2. Resolver
        3. Jumpgate
        4. LANclient
        5. DNSmaster
        6. Webbserver
        7. Mailserver
        8. DNSslave
        """
        while True:
            sc = int(input(str(sc_msg)))
            if sc == 1:
                proxy(paramlist[group_nr-1])
                break
            elif sc == 2:
                resolver(paramlist[group_nr-1])
                break
            elif sc == 3:
                jumpgate(paramlist[group_nr-1])
                break
            elif sc == 4:
                LANclient(paramlist[group_nr-1])
                break
            elif sc == 5:
                DNSmaster(paramlist[group_nr-1])
                break
            elif sc == 6:
                webbserver(paramlist[group_nr-1])
                break
            elif sc == 7:
                mailserver(paramlist[group_nr-1])
                break
            elif sc == 8:
                DNSslave(paramlist[group_nr-1])
                break
            else:
                print("Please choose a number between 1-8")
    else:
        deployAll(paramlist)
