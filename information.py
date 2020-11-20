import subprocess
import yaml
from fpdf import FPDF
import os
import argparse

"""
Script creating information tabels for every AWS environment
and an additional information sheet containing information tabels about every AWS environment for admin user
"""

if __name__ == "__main__":
    parser = argparse.ArgumentParser(prog='information', description='Get information from an AWS environment based on a parameter YAML-file')
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


    #create directory for the informationfiles if it doesn´t exists
    informationdir = "information"
    if os.path.isdir(informationdir):
        pass
    else:
        try:
            os.mkdir(informationdir)
        except Exception as e:
            print(e)

    servers = ["jumpgate", "mailserver", "webbserver", "DNSmaster", "DNSslave", "resolver", "LANclient"] #list of server-names that we want information about
    #loop through each AWS environment defined with parameters in the list paramlist and create an information table
    for params in paramlist:
        info = [] #list that will contain sublists with information (name, private Ipv4 address, public Ipv4 address and Ipv6 address) about every server
        LAN2 = [] #sublist to list "info" that will contain the secondary private Ipv4 and Ipv6 address of the LANclient
        LAN3 = [] #sublist to list "info" that will contain the third private Ipv4 and Ipv6 address of the LANclient
        LAN4 = [] #sublist to list "info" that will contain the forth private Ipv4 and Ipv6 address of the LANclient

        #loop through all servers in list "servers" and gather information about them in a list that will represent a row in the information table
        for server in servers:
            #get name, private Ipv4 address, public Ipv4 address and Ipv6 address of the current server in the loop
            server_info = subprocess.check_output("aws ec2 describe-instances --filter Name=tag:Name,Values="+params["GroupName"]+"-"+server+" Name=instance-state-code,Values=16 --query \"Reservations[].Instances[].[Tags[?Key=='Role']|[0].Value,NetworkInterfaces[].PrivateIpAddresses[].PrivateIpAddress,NetworkInterfaces[].PrivateIpAddresses[0].Association[].PublicIp,NetworkInterfaces[].Ipv6Addresses]\" --output text", shell=True)
            server_info = str(server_info) #convert output to string type
            server_info = server_info.strip("b\'\\r\\n") #strip string of unnecessary characters at the beginning and end of the string
            #split string into a list, containing name, private Ipv4 address, public Ipv4 address and Ipv6 address as seperate elements, at the separators "\t" and "\r\n"
            server_info = server_info.replace("\\t", "\\r\\n")
            server_info = server_info.split("\\r\\n")
            #devide the list of information about the LANclient into multiple lists that represents a row in the information table, since LANclient has multiple privte Ipv4 and Ipv6 addresses
            if server == "LANclient":
                for i in range(2, 6, 3):
                    LAN2.append("")
                    LAN2.append(server_info.pop(i)) #move an element (firs loop: private Ipv4 address, second loop: Ipv6 address) from the original LANclient list to the list contaning the second row of the LANclient
                for i in range(2, 5, 2):
                    LAN3.append("")
                    LAN3.append(server_info.pop(i)) #move an element (firs loop: private Ipv4 address, second loop: Ipv6 address) from the original LANclient list to the list contaning the third row of the LANclient
                for i in range(2, 4, 1):
                    LAN4.append("")
                    LAN4.append(server_info.pop(i)) #move an element (firs loop: private Ipv4 address, second loop: Ipv6 address) from the original LANclient list to the list contaning the fouth row of the LANclient
            #replace privte Ipv4 address od DNSslave with empty string
            if server == "DNSslave":
                server_info[1] = ""

            #get public Ipv4 address of proxy server
            PublicIp = subprocess.check_output("aws ec2 describe-instances --filter Name=tag:Name,Values="+params["GroupName"]+"-proxy"+" Name=instance-state-code,Values=16 --query \"Reservations[].Instances[].[NetworkInterfaces[].PrivateIpAddresses[0].Association[].PublicIp]\" --output text", shell=True)
            PublicIp = str(PublicIp) #convert output to string type
            PublicIp = PublicIp.strip("b\'\\r\\n") #strip string of unnecessary characters at the beginning and end of the string
            #insert public Ipv4 address of proxy server in the lists of information about the jumpgate, mailserver, webbserver and DNSmaster
            if server != "DNSslave" and server != "LANclient" and server != "resolver":
                server_info.insert(2, PublicIp)
            #insert empty string, in the place of public Ipv4 address, in the lists of information about the LANclient and resolver
            elif server == "LANclient" or server == "resolver":
                server_info.insert(2, "")
            info.append(server_info) #append the list of information about the current server in the loop
        #append the lists containg the second, third and fouth row of the LANclient information to the end of the "info" list
        info.append(LAN2)
        info.append(LAN3)
        info.append(LAN4)

        pdf = FPDF(orientation = "L") #create instance of FPDF class, landscape page format
        pdf.add_page() #add page for the information table of the current environment
        pdf.set_font("Arial", "B", 10)
        epw = pdf.w - 2*pdf.l_margin #effective page width, or just epw
        col_width = epw/4 #set column width to 1/4 of effective page width to distribute content evenly across table and page
        th = pdf.font_size #text height is the same as current font size
        #write the domain name at the top of the paper
        pdf.cell(col_width, th, txt="Domän", ln=0)
        pdf.set_font("Arial", "", 10)
        pdf.cell(col_width, th, txt=params["DomainName"], ln=1)
        pdf.set_font("Arial", "B", 10)
        pdf.ln(th)
        #write the heading of every coloumn in the table
        pdf.cell(col_width, th, txt="Server", ln=0)
        pdf.cell(col_width, th, txt="Privat Ipv4", ln=0)
        pdf.cell(col_width, th, txt="Publik Ipv4", ln=0)
        pdf.cell(col_width, th, txt="Ipv6", ln=1)
        pdf.set_font("Arial", "", 10)
        #white the information table
        for row in info:
            for element in row:
                pdf.cell(col_width, th, str(element), border=1) #write each element of the current sublist in list "info" in a row with borders drawn around the cell
            pdf.ln(th) #move to the beginning of next line
        pdf.output(informationdir+"/{}.pdf".format(str(params["GroupName"]))) #save information table in directory "information" and filename of the parameter "GroupName"


    pdf = FPDF(orientation = "L") #create instance of FPDF class, landscape page format
    pdf.add_page() #add page for the information sheet of the admin
    for params in paramlist:
        info = [] #list that will contain sublists with information (name, private Ipv4 address, public Ipv4 address and Ipv6 address) about every server
        LAN2 = [] #sublist to list "info" that will contain the secondary private Ipv4 and Ipv6 address of the LANclient
        LAN3 = [] #sublist to list "info" that will contain the third private Ipv4 and Ipv6 address of the LANclient
        LAN4 = [] #sublist to list "info" that will contain the forth private Ipv4 and Ipv6 address of the LANclient

        #loop through all servers in list "servers" and gather information about them in a list that will represent a row in the information table
        for server in servers:
            #get name, private Ipv4 address, public Ipv4 address and Ipv6 address of the current server in the loop
            server_info = subprocess.check_output("aws ec2 describe-instances --filter Name=tag:Name,Values="+params["GroupName"]+"-"+server+" Name=instance-state-code,Values=16 --query \"Reservations[].Instances[].[Tags[?Key=='Role']|[0].Value,NetworkInterfaces[].PrivateIpAddresses[].PrivateIpAddress,NetworkInterfaces[].PrivateIpAddresses[0].Association[].PublicIp,NetworkInterfaces[].Ipv6Addresses]\" --output text", shell=True)
            server_info = str(server_info) #convert output to string type
            server_info = server_info.strip("b\'\\r\\n") #strip string of unnecessary characters at the beginning and end of the string
            #split string into a list, containing name, private Ipv4 address, public Ipv4 address and Ipv6 address as seperate elements, at the separators "\t" and "\r\n"
            server_info = server_info.replace("\\t", "\\r\\n")
            server_info = server_info.split("\\r\\n")
            #devide the list of information about the LANclient into multiple lists that represents a row in the information table, since LANclient has multiple privte Ipv4 and Ipv6 addresses
            if server == "LANclient":
                for i in range(2, 6, 3):
                    LAN2.append("")
                    LAN2.append(server_info.pop(i)) #move an element (firs loop: private Ipv4 address, second loop: Ipv6 address) from the original LANclient list to the list contaning the second row of the LANclient
                for i in range(2, 5, 2):
                    LAN3.append("")
                    LAN3.append(server_info.pop(i)) #move an element (firs loop: private Ipv4 address, second loop: Ipv6 address) from the original LANclient list to the list contaning the third row of the LANclient
                for i in range(2, 4, 1):
                    LAN4.append("")
                    LAN4.append(server_info.pop(i)) #move an element (firs loop: private Ipv4 address, second loop: Ipv6 address) from the original LANclient list to the list contaning the fouth row of the LANclient
            #replace privte Ipv4 address od DNSslave with empty string
            if server == "DNSslave":
                server_info[1] = ""
            #get public Ipv4 address of proxy server
            PublicIp = subprocess.check_output("aws ec2 describe-instances --filter Name=tag:Name,Values="+params["GroupName"]+"-proxy"+" Name=instance-state-code,Values=16 --query \"Reservations[].Instances[].[NetworkInterfaces[].PrivateIpAddresses[0].Association[].PublicIp]\" --output text", shell=True)
            PublicIp = str(PublicIp) #convert output to string type
            PublicIp = PublicIp.strip("b\'\\r\\n") #strip string of unnecessary characters at the beginning and end of the string
            #insert public Ipv4 address of proxy server in the lists of information about the jumpgate, mailserver, webbserver and DNSmaster
            if server != "DNSslave" and server != "LANclient" and server != "resolver":
                server_info.insert(2, PublicIp)
            #insert empty string, in the place of public Ipv4 address, in the lists of information about the LANclient and resolver
            elif server == "LANclient" or server == "resolver":
                server_info.insert(2, "")
            info.append(server_info) #append the list of information about the current server in the loop
        #get public ILO Ipv4 address of proxy server
        ProxyILOIp = subprocess.check_output("aws ec2 describe-instances --filter Name=tag:Name,Values="+params["GroupName"]+"-proxy"+" Name=instance-state-code,Values=16 --query \"Reservations[].Instances[].[NetworkInterfaces[].PrivateIpAddresses[1].Association[].PublicIp]\" --output text", shell=True)
        ProxyILOIp = str(PublicIp) #convert output to string type
        ProxyILOIp = PublicIp.strip("b\'\\r\\n") #strip string of unnecessary characters at the beginning and end of the string
        ProxyILO = ["Proxy ILO", "", ProxyILOIp, ""] #create list for the last row in the informationtable, containing the ILO address of the proxy server
        #append the lists containg the second, third and fouth row of the LANclient information, and the list containing the ILO address of the proxy server to the end of the "info" list
        info.append(LAN2)
        info.append(LAN3)
        info.append(LAN4)
        info.append(ProxyILO)

        pdf.set_font("Arial", "B", 10)
        epw = pdf.w - 2*pdf.l_margin #effective page width, or just epw
        col_width = epw/4 #set column width to 1/4 of effective page width to distribute content evenly across table and page
        th = pdf.font_size #text height is the same as current font size
        #write the domain name at the top of the paper
        pdf.cell(col_width, th, txt="Domän", ln=0)
        pdf.set_font("Arial", "", 10)
        pdf.cell(col_width, th, txt=params["DomainName"], ln=1)
        pdf.set_font("Arial", "B", 10)
        #write the key name under the domain name
        pdf.cell(col_width, th, txt="Nyckel", ln=0)
        pdf.set_font("Arial", "", 10)
        pdf.cell(col_width, th, txt=params["StudentKeyName"], ln=1)
        pdf.set_font("Arial", "B", 10)
        pdf.ln(th)
        #write the heading of every coloumn in the table
        pdf.cell(col_width, th, txt="Server", ln=0)
        pdf.cell(col_width, th, txt="Privat Ipv4", ln=0)
        pdf.cell(col_width, th, txt="Publik Ipv4", ln=0)
        pdf.cell(col_width, th, txt="Ipv6", ln=1)
        pdf.set_font("Arial", "", 10)
        #write the information table of each group
        for row in info:
            for object in row:
                pdf.cell(col_width, th, str(object), border=1)#write each element of the current sublist in list "info" in a row with borders drawn around the cell
            pdf.ln(th) #move to the beginning of next line
        pdf.ln(4*th) #line break equivalent to 4 lines, to make space between each environment table
    pdf.output(informationdir+"/admin.pdf") #save all information tables in a single information sheet in directory "information" and filename "admin"
