#!/usr/bin/env python3
import yaml
import subprocess
import argparse
import time


def VPC(params):
    StackName = params["GroupName"] + "-VPC"
    print("Deleting "+StackName)
    Subnet_As_ID = subprocess.check_output("aws ec2 describe-subnets --filter Name=tag:Name,Values="+StackName+" --query Subnets[].Ipv6CidrBlockAssociationSet[].AssociationId --output text", shell=True)
    Subnet_As_ID = str(Subnet_As_ID)
    Subnet_As_ID = Subnet_As_ID.strip("b\'\\r\\n")
    Subnet_As_ID = Subnet_As_ID.split("\\t")
    for ID in Subnet_As_ID:
        print(ID)
        subprocess.run(["aws", "ec2", "disassociate-subnet-cidr-block", "--association-id", ID])
    VPC_As_ID = subprocess.check_output("aws ec2 describe-vpcs --filter Name=tag:Name,Values="+StackName+" --query Vpcs[].Ipv6CidrBlockAssociationSet[].AssociationId --output text", shell=True)
    VPC_As_ID = str(VPC_As_ID)
    VPC_As_ID = VPC_As_ID.strip("b\'\\r\\n")
    subprocess.run(["aws", "ec2", "disassociate-vpc-cidr-block", "--association-id", VPC_As_ID])
    Pu_CB_state = ""
    LAN_CB_state = ""
    DMZ_CB_state = ""
    VPC_CB_state = ""
    while Pu_CB_state != "disassociated" and LAN_CB_state !="disassociated" and DMZ_CB_state != "disassociated" and VPC_CB_state != "disassociated":
        Pu_CB_state = subprocess.check_output("aws ec2 describe-subnets --filter Name=tag:Name,Values="+StackName+" --query Subnets[].Ipv6CidrBlockAssociationSet[].Ipv6CidrBlockState.State --output text", shell=True)
        Pu_CB_state = str(Pu_CB_state)
        Pu_CB_state = Pu_CB_state.strip("b\'\\r\\n")
        LAN_CB_state = subprocess.check_output("aws ec2 describe-vpcs --filter Name=tag:Name,Values="+StackName+" --query Subnets[].Ipv6CidrBlockAssociationSet[].Ipv6CidrBlockState.State --output text", shell=True)
        LAN_CB_state = str(LAN_CB_state)
        LAN_CB_state = LAN_CB_state.strip("b\'\\r\\n")
        DMZ_CB_state = subprocess.check_output("aws ec2 describe-vpcs --filter Name=tag:Name,Values="+StackName+" --query Subnets[].Ipv6CidrBlockAssociationSet[].Ipv6CidrBlockState.State --output text", shell=True)
        DMZ_CB_state = str(DMZ_CB_state)
        DMZ_CB_state = DMZ_CB_state.strip("b\'\\r\\n")
        VPC_CB_state = subprocess.check_output("aws ec2 describe-vpcs --filter Name=tag:Name,Values="+StackName+" --query Vpcs[].Ipv6CidrBlockAssociationSet[].Ipv6CidrBlockState.State --output text", shell=True)
        VPC_CB_state = str(VPC_CB_state)
        VPC_CB_state = VPC_CB_state.strip("b\'\\r\\n")
        time.sleep(5)
    args = ["aws", "cloudformation", "delete-stack", "--stack-name", StackName]
    subprocess.run(args)

def VPCSlaveDNS():
    StackName = "VPCSlaveDNS"
    print("Deleting "+StackName)
    Subnet_As_ID = subprocess.check_output("aws ec2 describe-subnets --filter Name=tag:Name,Values="+StackName+" --query Subnets[].Ipv6CidrBlockAssociationSet[].AssociationId --output text", shell=True)
    Subnet_As_ID = str(Subnet_As_ID)
    Subnet_As_ID = Subnet_As_ID.strip("b\'\\r\\n")
    Subnet_As_ID = Subnet_As_ID.split("\\t")
    for ID in Subnet_As_ID:
        print(ID)
        subprocess.run(["aws", "ec2", "disassociate-subnet-cidr-block", "--association-id", ID])
    VPC_As_ID = subprocess.check_output("aws ec2 describe-vpcs --filters Name=tag:Name,Values="+StackName+" --query Vpcs[].Ipv6CidrBlockAssociationSet[].AssociationId --output text", shell=True)
    VPC_As_ID = str(VPC_As_ID)
    VPC_As_ID = VPC_As_ID.strip("b\'\\r\\n")
    subprocess.run(["aws", "ec2", "disassociate-vpc-cidr-block", "--association-id", VPC_As_ID])
    Pu_CB_state = ""
    LAN_CB_state = ""
    DMZ_CB_state = ""
    VPC_CB_state = ""
    while Pu_CB_state != "disassociated" and LAN_CB_state !="disassociated" and DMZ_CB_state != "disassociated" and VPC_CB_state != "disassociated":
        Pu_CB_state = subprocess.check_output("aws ec2 describe-subnets --filter Name=tag:Name,Values="+StackName+" --query Subnets[].Ipv6CidrBlockAssociationSet[].Ipv6CidrBlockState.State --output text", shell=True)
        Pu_CB_state = str(Pu_CB_state)
        Pu_CB_state = Pu_CB_state.strip("b\'\\r\\n")
        LAN_CB_state = subprocess.check_output("aws ec2 describe-vpcs --filter Name=tag:Name,Values="+StackName+" --query Subnets[].Ipv6CidrBlockAssociationSet[].Ipv6CidrBlockState.State --output text", shell=True)
        LAN_CB_state = str(LAN_CB_state)
        LAN_CB_state = LAN_CB_state.strip("b\'\\r\\n")
        DMZ_CB_state = subprocess.check_output("aws ec2 describe-vpcs --filter Name=tag:Name,Values="+StackName+" --query Subnets[].Ipv6CidrBlockAssociationSet[].Ipv6CidrBlockState.State --output text", shell=True)
        DMZ_CB_state = str(DMZ_CB_state)
        DMZ_CB_state = DMZ_CB_state.strip("b\'\\r\\n")
        VPC_CB_state = subprocess.check_output("aws ec2 describe-vpcs --filter Name=tag:Name,Values="+StackName+" --query Vpcs[].Ipv6CidrBlockAssociationSet[].Ipv6CidrBlockState.State --output text", shell=True)
        VPC_CB_state = str(VPC_CB_state)
        VPC_CB_state = VPC_CB_state.strip("b\'\\r\\n")
        time.sleep(5)
    args = ["aws", "cloudformation", "delete-stack", "--stack-name", StackName]
    subprocess.run(args)

def securityGroups(params):
    StackName = params["GroupName"] + "-SG"
    print("Deleting "+StackName)
    args = ["aws", "cloudformation", "delete-stack", "--stack-name", StackName]
    subprocess.run(args)

def proxy(params):
    StackName = params["GroupName"] + "-proxy"
    print("Deleting "+StackName)
    args = ["aws", "cloudformation", "delete-stack", "--stack-name", StackName]
    subprocess.run(args)

def resolver(params):
    StackName = params["GroupName"] + "-resolver"
    print("Deleting "+StackName)
    args = ["aws", "cloudformation", "delete-stack", "--stack-name", StackName]
    subprocess.run(args)

def jumpgate(params):
    StackName = params["GroupName"] + "-jumpgate"
    print("Deleting "+StackName)
    args = ["aws", "cloudformation", "delete-stack", "--stack-name", StackName]
    subprocess.run(args)

def LANclient(params):
    StackName = params["GroupName"] + "-LANclient"
    print("Deleting "+StackName)
    args = ["aws", "cloudformation", "delete-stack", "--stack-name", StackName]
    subprocess.run(args)

def DNSmaster(params):
    StackName = params["GroupName"] + "-DNSmaster"
    print("Deleting "+StackName)
    args = ["aws", "cloudformation", "delete-stack", "--stack-name", StackName]
    subprocess.run(args)

def webbserver(params):
    StackName = params["GroupName"] + "-webbserver"
    print("Deleting "+StackName)
    args = ["aws", "cloudformation", "delete-stack", "--stack-name", StackName]
    subprocess.run(args)

def mailserver(params):
    StackName = params["GroupName"] + "-mailserver"
    print("Deleting "+StackName)
    args = ["aws", "cloudformation", "delete-stack", "--stack-name", StackName]
    subprocess.run(args)

def DNSslave(params):
    StackName = params["GroupName"] + "-DNSslave"
    print("Deleting "+StackName)
    args = ["aws", "cloudformation", "delete-stack", "--stack-name", StackName]
    subprocess.run(args)

def deleteAll(paramlist):
    for params in paramlist:
        DNSslave(params)
        mailserver(params)
        webbserver(params)
        DNSmaster(params)
        LANclient(params)
        jumpgate(params)
        resolver(params)
        proxy(params)
        subprocess.run(["aws", "cloudformation", "wait", "stack-delete-complete", "--stack-name", params["GroupName"]+"-DNSslave"])
        subprocess.run(["aws", "cloudformation", "wait", "stack-delete-complete", "--stack-name", params["GroupName"]+"-mailserver"])
        subprocess.run(["aws", "cloudformation", "wait", "stack-delete-complete", "--stack-name", params["GroupName"]+"-webbserver"])
        subprocess.run(["aws", "cloudformation", "wait", "stack-delete-complete", "--stack-name", params["GroupName"]+"-DNSmaster"])
        subprocess.run(["aws", "cloudformation", "wait", "stack-delete-complete", "--stack-name", params["GroupName"]+"-LANclient"])
        subprocess.run(["aws", "cloudformation", "wait", "stack-delete-complete", "--stack-name", params["GroupName"]+"-jumpgate"])
        subprocess.run(["aws", "cloudformation", "wait", "stack-delete-complete", "--stack-name", params["GroupName"]+"-resolver"])
        subprocess.run(["aws", "cloudformation", "wait", "stack-delete-complete", "--stack-name", params["GroupName"]+"-proxy"])
        securityGroups(params)
        subprocess.run(["aws", "cloudformation", "wait", "stack-delete-complete", "--stack-name", params["GroupName"]+"-SG"])
        VPC(params)
        subprocess.run(["aws", "cloudformation", "wait", "stack-delete-complete", "--stack-name", params["GroupName"]+"-VPC"])
    VPCSlaveDNS()
    subprocess.run(["aws", "cloudformation", "wait", "stack-delete-complete", "--stack-name", "VPCSlaveDNS"])

def deleteOne(params):
    DNSslave(params)
    mailserver(params)
    webbserver(params)
    DNSmaster(params)
    LANclient(params)
    jumpgate(params)
    resolver(params)
    proxy(params)
    subprocess.run(["aws", "cloudformation", "wait", "stack-delete-complete", "--stack-name", params["GroupName"]+"-DNSslave"])
    subprocess.run(["aws", "cloudformation", "wait", "stack-delete-complete", "--stack-name", params["GroupName"]+"-mailserver"])
    subprocess.run(["aws", "cloudformation", "wait", "stack-delete-complete", "--stack-name", params["GroupName"]+"-webbserver"])
    subprocess.run(["aws", "cloudformation", "wait", "stack-delete-complete", "--stack-name", params["GroupName"]+"-DNSmaster"])
    subprocess.run(["aws", "cloudformation", "wait", "stack-delete-complete", "--stack-name", params["GroupName"]+"-LANclient"])
    subprocess.run(["aws", "cloudformation", "wait", "stack-delete-complete", "--stack-name", params["GroupName"]+"-jumpgate"])
    subprocess.run(["aws", "cloudformation", "wait", "stack-delete-complete", "--stack-name", params["GroupName"]+"-resolver"])
    subprocess.run(["aws", "cloudformation", "wait", "stack-delete-complete", "--stack-name", params["GroupName"]+"-proxy"])
    securityGroups(params)
    subprocess.run(["aws", "cloudformation", "wait", "stack-delete-complete", "--stack-name", params["GroupName"]+"-SG"])
    VPC(params)
    subprocess.run(["aws", "cloudformation", "wait", "stack-delete-complete", "--stack-name", params["GroupName"]+"-VPC"])


if __name__ == "__main__":
    parser = argparse.ArgumentParser(prog='delete', description='Delete an AWS environment from YAML-file')
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
    Select an option for deletion (default: 1)
    1. Delete all environments
    2. Delete a specific groups environment
    3. Delete a specific server for a specific group
    """
    choice = int(input(str(choice_msg)) or 1)
    if choice == 1:
        deleteAll(paramlist)
    elif choice == 2:
        group_nr = int(input("Which group?"))
        deleteOne(paramlist[group_nr-1])
    elif choice == 3:
        group_nr = int(input("Which group?"))
        sc_msg = """
        Select an option for which server to delete
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
        deleteAll(paramlist)
