import subprocess
import requests
import hashlib
import os
import time
import urllib

class P2PClient:
    """
    P2PClient is a class that provides methods to interact with a P2P file transfer API.
    It supports creating requests, uploading files via Aspera or HTTP, downloading files,
    and querying request details.

    Methods:
        __init__: Initializes the P2PClient with the given endpoint and SSL verification flag.
        create_request: Creates a P2P request using the provided access token and request data.
        calculate_sha1_hash: Calculates the SHA1 hash of the specified file.
        upload_files_by_aspera: Uploads files using Aspera with the given details.
        upload_files_by_http: Uploads files via HTTP with the given details.
        upload_file_in_chunks: Uploads a file in chunks via HTTP.
        get_request_detail_by_request_id: Retrieves the details of a request by request ID.
        check_upload_completion: Checks the completion status of an upload request.
        query_my_upload_requests: Queries upload requests based on filter criteria.
        query_my_received_requests: Queries received requests based on filter criteria.
        get_aspera_download_token_by_request_id: Retrieves Aspera download tokens by request ID.
        download_files_by_aspera: Downloads files via Aspera by request ID.
        get_http_download_token_by_request_id: Retrieves HTTP download tokens by request ID.
        download_files_by_http: Downloads files via HTTP by request ID.
    """

    def __init__(self, endpoint, ssl_verify):
        """
        Initializes the P2PClient with the given endpoint.

        Args:
            endpoint (str): The endpoint URL for the api server.
            ssl_verify: Https ssl_verify flag. False for disable ssl verify, True for enable ssl verify.
        """
        self.endpoint = endpoint
        self.ssl_verify = ssl_verify

    def create_request(self, access_token, request):
        """
        Create a P2P request using the provided access token and request data.

        Args:
            access_token (str): The access token for authentication.
            request (dict): The request data. Example:
                {
                    "FileHashAlgorithm": "",  # The algorithm used for file hashing (e.g., "sha1sum").
                    "Files": [  # List of files to be included in the request.
                        {
                            "FileName": "123.txt",  # The name of the file.
                            "Hash": "d58411c9c36c679a80c6cdff748fbe7efa2a0205",  # The hash of the file.
                            "Size": 6402951  # The size of the file in bytes.
                        }
                    ],
                    "Mail": {  # Email details for the request.
                        "BccList": [],  # List of BCC email addresses.
                        "Body": "Demo create request",  # The body of the email.
                        "CcList": [],  # List of CC email addresses.
                        "Title": "Demo create request",  # The title of the email.
                        "ToList": [  # List of recipient email addresses.
                            {
                                "Email": "liang.wang@mediatek.com"  # The email address of the recipient.
                            }
                        ]
                    },
                    "Subject": "demo",  # The subject of the request.
                    "System": "API"  # The system from which the request is made.
                }

        Returns:
            dict: The response JSON. Example:
                {
                    "code": 200,
                    "count": 1,
                    "message": "Success",
                    "results": {
                        "RequestId": "REQ20000087953", # request id
                        "ascpargs": "{{ascp.exe path}} --ignore-host-key -d -T -U 2 -k 2 --mode=send -i {{openssh file path}} -P 33001 -O 33001  --user=MFT3_SYS --host=172.21.144.160 -W ATM3_ACsJF20tTCTs_2TO9PT42bWFA6gloTs_04PexohumM4gG8AAFJdqgQ6mw_NH4bArS5Vikb_3MTA -l 100000 -m 100000 \"{file dir}/123.txt\" \"/aspera/DATA/MFT3_SYS/REQ20000087953\"", # aspera command line template
                        "authentication": "token",
                        "capSpeed": 100000,
                        "cipher": "aes-128",
                        "create_dir": "true",
                        "destination_root": "/aspera/DATA/MFT3_SYS/REQ20000087953",
                        "direction": "send",
                        "fasp_port": 33001,
                        "http_fallback": false,
                        "lock_min_rate": true,
                        "minSpeed": 100000,
                        "min_rate_kbps": 0,
                        "paths": [
                            {}
                        ],
                        "rate_policy": "fair",
                        "rate_policy_allowed": "fixed",
                        "remote_host": "xxx",
                        "remote_user": "xxx",
                        "resume": "sparse_checksum",
                        "source_root": "",
                        "speed": 100000,
                        "ssh_port": 33001,
                        "sshfp": null,
                        "tags": null,
                        "target_rate_kbps": 10000,
                        "token": "ATM3_ACsJF20tTCTs_2TO9PT42bWFA6gloTs_04PexohumM4gG8AAFJdqgQ6mw_NH4bArS5Vikb_3MTA"
                    }
                }
        """
        url = self.endpoint + "/Customer/FEX/P2P/v1/sharedRequests"
        headers = {
            "Content-Type": "application/json; charset=utf-8",
            "Authorization": f"Bearer {access_token}"
        }
        response = requests.post(url, json=request, headers=headers, verify=self.ssl_verify)
        return response.json()

    def calculate_sha1_hash(self, filename):
        """
        Calculates the hash of the specified file using the given algorithm.

        Args:
            filename (str): The name of the file to calculate the hash for.
            algorithm (str): The hash algorithm to use (sha1').

        Returns:
            str: The hash of the file.
        """
        hash_func = hashlib.sha1()
        with open(filename, 'rb') as f:
            for chunk in iter(lambda: f.read(4096), b""):
                hash_func.update(chunk)
        return hash_func.hexdigest()

    def upload_files_by_aspera(self, access_token, aspera_install_dir, subject, body, toList, ccList, bccList, files, verifyBySha1, 
                               retry_delay=120, retry_attempts=3):
        """
        Uploads files using Aspera with the given details.

        Args:
            subject (str): The subject of the email.
            body (str): The body of the email.
            toList (list): The list of recipients.
            ccList (list): The list of CC recipients.
            bccList (list): The list of BCC recipients.
            files (list): The list of files to upload.
            verifyBySha1 (bool): Whether to verify the files by SHA1.

        Returns:
            dict: The request data.
        """

        if not subject:
            raise ValueError("Subject is required")
        if not toList:
            raise ValueError("ToList is required")
        if not isinstance(files, list) or not files:
            raise ValueError("Files are required and must be a non-empty list")

        # Check if required Aspera files exist
        ascp_path = os.path.join(aspera_install_dir, 'bin', 'ascp.exe')
        openssh_path = os.path.join(aspera_install_dir, 'etc', 'asperaweb_id_dsa.openssh')
        if not os.path.exists(ascp_path) or not os.path.exists(openssh_path):
            raise FileNotFoundError("Required Aspera files are missing")

        request = {
            "FileHashAlgorithm": "sha1sum" if verifyBySha1 else "",
            "Files": [],
            "Mail": {
                "BccList": bccList,
                "Body": body, 
                "CcList": ccList,
                "Title": subject, 
                "ToList": toList
            },
            "Subject": subject,
            "System": "API"
        }

        for filename in files:
            file_info = {
                "FileName": os.path.basename(filename),
                "Size": os.path.getsize(filename)
            }
            if verifyBySha1:
                file_info["Hash"] = self.calculate_sha1_hash(filename)
            request["Files"].append(file_info)
        
        create_req_response = self.create_request(access_token, request)
        if 'results' not in create_req_response:
            raise ValueError("Failed to create request")
        aspera_upload_param = create_req_response["results"]
        
        if "ascpargs" not in aspera_upload_param:
            raise ValueError("ascpargs is missing in the response")
        ascp_temp = aspera_upload_param["ascpargs"]
        cmd = ascp_temp.replace("{{ascp.exe path}}", ascp_path).replace("{{openssh file path}}", openssh_path)
        for filename in files:
            cmd = cmd.replace(f"{{file dir}}/{os.path.basename(filename)}", filename)
        
        process = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        stdout, stderr = process.communicate()

        # Retry logic
        for attempt in range(1, retry_attempts + 1):
            process = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
            stdout, stderr = process.communicate()

            if process.returncode == 0:
                return aspera_upload_param["RequestId"], stdout.decode('utf-8'), ''
            
            if attempt >= retry_attempts:
                return aspera_upload_param["RequestId"], stdout.decode('utf-8'), stderr.decode('utf-8')
            
            time.sleep(retry_delay)

    def upload_files_by_http(self, access_token, subject, body, toList, ccList, bccList, files, verifyBySha1, 
                             chunk_size=1024 * 1024 * 100, chunk_retry_attempts=100):
        """
        Upload files via HTTP.

        Parameters:
        access_token (str): Access token for authorization.
        subject (str): Subject of the email.
        body (str): Body of the email.
        toList (list): List of recipient email addresses.
        ccList (list): List of CC email addresses.
        bccList (list): List of BCC email addresses.
        files (list): List of file paths to be uploaded.
        verifyBySha1 (bool): Whether to verify files by SHA1 hash.
        chunk_size (int, optional): Size of each chunk in bytes. Defaults to 100MB.
        chunk_retry_attempts (int, optional): Number of retry attempts if the upload of a chunk fails. Defaults to 10.

        Returns:
        tuple: A tuple containing the request ID and a list of errors that occurred during file upload.

        Raises:
        ValueError: If required parameters are missing or if the request creation or file upload fails.
        """
        if not subject:
            raise ValueError("Subject is required")
        if not toList:
            raise ValueError("ToList is required")
        if not isinstance(files, list) or not files:
            raise ValueError("Files are required and must be a non-empty list")
        
        request = {
            "FileHashAlgorithm": "sha1sum" if verifyBySha1 else "",
            "Files": [],
            "Mail": {
                "BccList": bccList,
                "Body": body, 
                "CcList": ccList,
                "Title": subject, 
                "ToList": toList
            },
            "Subject": subject,
            "System": "API"
        }

        for filename in files:
            file_size = os.path.getsize(filename)
            file_info = {
            "FileName": os.path.basename(filename),
            "Size": file_size,
            "TransferType": 2
            }
            if verifyBySha1:
                file_info["Hash"] = self.calculate_sha1_hash(filename)
            request["Files"].append(file_info)
        
        create_req_response = self.create_request(access_token, request)
        if 'results' not in create_req_response:
            raise ValueError("Failed to create request")
        transfer_param = create_req_response["results"]
        if 'transferServerParameter' not in transfer_param:
            raise ValueError("transferServerParameters is missing in the response")
        transfer_server_param = transfer_param["transferServerParameter"]

        error_list = []

        for file in files:
            try:
                self.upload_file_in_chunks(transfer_server_param["url"], transfer_server_param["token"], transfer_server_param["requestId"], file, chunk_size, chunk_retry_attempts)
            except Exception as e:
                error_list.append({"file": file, "error": str(e)})

        if error_list:
            raise ValueError(f"Errors occurred during file upload: {error_list}")

        return transfer_param["RequestId"], error_list

    def upload_file_in_chunks(self, url, token, requestId, filename, chunk_size, retry_attempts):
        """
        Upload a file in chunks.

        Parameters:
        url (str): The URL to which the file chunks will be uploaded.
        token (str): Authorization token for the upload.
        requestId (str): Request ID to identify the upload request.
        filename (str): The path to the file to be uploaded.
        chunk_size (int): The size of each chunk in bytes.
        retry_attempts (int): The number of retry attempts if the upload of a chunk fails.

        Upload Parameters:
        - requestId (str): The request ID.
        - fileName (str): The name of the file being uploaded.
        - chunk (int): The current chunk index.
        - chunks (int): The total number of chunks.
        - size (int): The total size of the file in bytes.
        - chunkSize (int): The size of the current chunk in bytes.

        Raises:
        ValueError: If a chunk fails to upload after the specified number of retry attempts.
        """
        headers = {
            "Authorization": f"Token {token}"
        }
        
        file_size = os.path.getsize(filename)
        file_name = os.path.basename(filename)
        total_chunks = (file_size + chunk_size - 1) // chunk_size
        
        with open(filename, 'rb') as f:
            for chunk_index in range(total_chunks):
                chunk_data = f.read(chunk_size)
                files = {
                    'file': (file_name, chunk_data)
                }
                data = {
                    'requestId': requestId,
                    'fileName': file_name,
                    'chunk': chunk_index,
                    'chunks': total_chunks,
                    'size': file_size,
                    'chunkSize': chunk_size
                }
                
                for _ in range(1, retry_attempts + 1):  # 重试10次
                    response = requests.post(url, headers=headers, files=files, data=data)
                    if response.status_code != 200:
                        time.sleep(1)
                    break;                    
                else:
                    raise ValueError(f"Failed to upload chunk {chunk_index + 1} of file {file_name} after 10 attempts: {response.text}")

    def get_request_detail_by_request_id(self, access_token, request_id):
        """
        Get the details of a request by request ID.

        Parameters:
        access_token (str): Access token for authorization.
        request_id (str): Request ID to identify the request.

        Returns:
        dict: A dictionary containing the details of the request.

        Raises:
        ValueError: If the response JSON does not contain 'results' or if the request fails.
        """
        encoded_request_id = urllib.parse.quote(request_id)
        url = self.endpoint + "/Customer/FEX/P2P/v1/requests/" + encoded_request_id
        headers = {
            "Content-Type": "application/json; charset=utf-8",
            "Authorization": f"Bearer {access_token}"
        }

        response = requests.get(url, headers=headers, verify=self.ssl_verify, timeout=(60, 300))
        
        if response.status_code == 200:
            response_data = response.json()
            if 'results' in response_data:
                return response_data['results']
            else:
                raise ValueError(f"Response JSON does not contain 'results': {response.text}")
        else:
            raise ValueError(f"Failed to get request details, status code: {response.status_code}, response: {response.text}")

    def check_upload_completion(self, access_token_provider, request_id, wait_for_completion=False, check_upload_completion_interval=60, check_upload_completion_attempts=10):
        """
        Check the completion status of an upload request.

        Parameters:
        access_token_provider (callable): A callable that returns the access token for authorization.
        request_id (str): Request ID to identify the upload request.
        wait_for_completion (bool, optional): Whether to wait for the upload to complete. Defaults to False.
        check_upload_completion_interval (int, optional): Interval in seconds between status checks. Defaults to 60.
        check_upload_completion_attempts (int, optional): Number of attempts to check the upload status. Defaults to 10.

        Returns:
        bool: True if the upload is complete and the status is "Ready", False otherwise.

        Raises:
        ValueError: If the request status is "Error".
        """
        access_token = access_token_provider()
        request_detail = self.get_request_detail_by_request_id(access_token, request_id)
        if 'Status' in request_detail:
            if request_detail['Status'] == "Error":
                return False
            if wait_for_completion:
                attempt = 1
                while request_detail['Status'] != "Ready" and attempt <= check_upload_completion_attempts:
                    time.sleep(check_upload_completion_interval)
                    access_token = access_token_provider()
                    request_detail = self.get_request_detail_by_request_id(access_token, request_id)
                    attempt += 1
            return request_detail['Status'] == "Ready"
        return False

    def query_my_upload_requests(self, access_token, filter):
        """
        Query upload requests.

        Parameters:
        access_token (str): Access token for authorization.
        filter (dict): Filter criteria for querying upload requests. The filter can include the following keys:
            - PageIndex (int): The index of the page to retrieve.
            - PageSize (int): The number of items per page.
            - Status (str): The status of the requests to filter by.
            - RequestIn (str): The request ID to filter by, semicolon-separated list of request IDs.
            - RecceiveLike (str): The receiver's name or email to filter by.
            - SubjectLike (str): The subject to filter by.
            - ShareTimeFrom (str): The start time to filter by (ISO 8601 format).
            - ShareTimeTo (str): The end time to filter by (ISO 8601 format).
            - FileNameLike (str): The file name to filter by.

        Returns:
        list: A list of dictionaries containing the results of the query.
        {
            "Count": int,  # The total number of upload requests.
            "DataList": [
                {
                    "RequestId": str,  # The unique identifier for the request.
                    "Subject": str,  # The subject of the request.
                    "ShareTime": str,  # The time when the request was shared (YYYY/MM/DD HH:MM:SS).
                    "ExpireTime": str or None,  # The time when the request will expire (YYYY/MM/DD HH:MM:SS) or None if not applicable.
                    "Status": str,  # The status of the request (e.g., "Uploading").
                    "System": str,  # The system from which the request was made.
                    "Count": int  # The total count of requests.
                }
            ]
        }
        Raises:
        ValueError: If the response JSON does not contain 'results' or if the request fails.
        """
        url = self.endpoint + "/Customer/FEX/P2P/v1/sharedRequestsQuerier"
        headers = {
            "Content-Type": "application/json; charset=utf-8",
            "Authorization": f"Bearer {access_token}"
        }

        response = requests.post(url, headers=headers, json=filter, verify=self.ssl_verify)
        
        if response.status_code == 200:
            response_data = response.json()
            if 'results' in response_data:
                return response_data['results']
            else:
                raise ValueError(f"Response JSON does not contain 'results': {response.text}")
        else:
            raise ValueError(f"Failed to get request details, status code: {response.status_code}, response: {response.text}")
    
    def query_my_received_requests(self, access_token, filter):
        """
        Query received requests.

        Parameters:
        access_token (str): Access token for authorization.
        filter (dict): Filter criteria for querying received requests. The filter can include the following keys:
            - PageIndex (int): The index of the page to retrieve.
            - PageSize (int): The number of items per page.
            - Status (str): The status of the requests to filter by.
            - RequestIn (str): The request ID to filter by.
            - SendLike (str): The sender's name or email to filter by.
            - SubjectLike (str): The subject to filter by.
            - ReceivedTimeFrom (str): The start time to filter by (ISO 8601 format).
            - ReceivedTimeTo (str): The end time to filter by (ISO 8601 format).
            - FileNameLike (str): The file name to filter by.

        Returns:
        list: A list of dictionaries containing the results of the query.
             {
                "Count": int,  # The total number of received requests.
                "DataList": [
                    {
                        "MailId": str,  # The unique identifier for the mail.
                        "UserEmail": str,  # The email address of the user who received the request.
                        "System": str,  # The system from which the request was made.
                        "SendTime": str,  # The time when the request was sent (YYYY/MM/DD HH:MM:SS).
                        "ExpireTime": str,  # The time when the request will expire (YYYY/MM/DD HH:MM:SS).
                        "Title": str,  # The title of the request.
                        "Status": str,  # The status of the request (e.g., "Ready").
                        "RequestId": str,  # The unique identifier for the request.
                        "Downloaded": bool,  # Indicates whether the request has been downloaded.
                        "Count": int  # The total count of requests.
                    }
                ]
            }
        Raises:
        ValueError: If the response JSON does not contain 'results' or if the request fails.
        """
        url = self.endpoint + "/Customer/FEX/P2P/v1/receivedRequestsQuerier"
        headers = {
            "Content-Type": "application/json; charset=utf-8",
            "Authorization": f"Bearer {access_token}"
        }

        response = requests.post(url, headers=headers, json=filter, verify=self.ssl_verify)
        
        if response.status_code == 200:
            response_data = response.json()
            if 'results' in response_data:
                return response_data['results']
            else:
                raise ValueError(f"Response JSON does not contain 'results': {response.text}")
        else:
            raise ValueError(f"Failed to get request details, status code: {response.status_code}, response: {response.text}") 
    
    def get_aspera_download_token_by_request_id(self, access_token, request_id, filenames=None):
        """
        Get Aspera download token by request ID.

        Parameters:
        access_token (str): Access token for authorization.
        request_id (str): Request ID to identify the download request.
        filenames (list, optional): List of filenames to specify the files to download. Defaults to None.

        Returns:
        list: A list of dictionaries containing download tokens.
            - "ascpargs" (str): The Aspera command arguments for downloading the file.
            - "CompleteDate" (str): The date and time when the download was completed (YYYY/MM/DD HH:MM:SS).

        Raises:
        ValueError: If the response JSON does not contain 'results' or if the request fails.
        """
        encoded_request_id = urllib.parse.quote(request_id)
        url = self.endpoint + f"/Customer/FEX/P2P/v1/requests/{encoded_request_id}/AsperaDownloadCommands"
        headers = {
            "Content-Type": "application/json; charset=utf-8",
            "Authorization": f"Bearer {access_token}"
        }
        if filenames:
            headers["filenames"] = "/".join(filenames)
        
        response = requests.get(url, headers=headers, verify=self.ssl_verify)
        if response.status_code == 200:
            response_data = response.json()
            if 'results' in response_data:
                return response_data['results']
            else:
                raise ValueError(f"Response JSON does not contain 'results': {response.text}")
        else:
            raise ValueError(f"Failed to get request details, status code: {response.status_code}, response: {response.text}")

    def download_files_by_aspera(self, access_token, request_id, filenames, aspera_install_dir, target_dir, retry_delay=120, retry_attempts=3):
        """
        Download files via Aspera by request ID.

        Parameters:
        access_token (str): Access token for authorization.
        request_id (str): Request ID to identify the download request.
        filenames (list): List of filenames to specify the files to download.
        aspera_install_dir (str): Directory where Aspera is installed.
        target_dir (str): Directory where the downloaded files will be saved.
        retry_delay (int, optional): Delay in seconds between retry attempts. Defaults to 120.
        retry_attempts (int, optional): Number of retry attempts if the download fails. Defaults to 3.

        Returns:
        tuple: A tuple containing the standard output and standard error of the Aspera command.

        Raises:
        FileNotFoundError: If the required Aspera files are missing.
        ValueError: If 'ascpargs' is missing in the response.
        """
        # Check if target directory exists, if not, create it
        if not os.path.exists(target_dir):
            os.makedirs(target_dir)

        # Check if required Aspera files exist
        ascp_path = os.path.join(aspera_install_dir, 'bin', 'ascp.exe')
        openssh_path = os.path.join(aspera_install_dir, 'etc', 'asperaweb_id_dsa.openssh')
        if not os.path.exists(ascp_path) or not os.path.exists(openssh_path):
            raise FileNotFoundError("Required Aspera files are missing")

        download_param = self.get_aspera_download_token_by_request_id(access_token, request_id, filenames)
        if 'ascpargs' not in download_param:
            raise ValueError("ascpargs is missing in the response")
        ascp_temp = download_param["ascpargs"]
        cmd = ascp_temp.replace("{ascp.exe path}", ascp_path).replace("{openssh file path}", openssh_path).replace("{download dir}", f"\"{target_dir}\"")

        # Retry logic
        for attempt in range(1, retry_attempts + 1):
            process = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
            stdout, stderr = process.communicate()

            if process.returncode == 0:
                return stdout.decode('utf-8'), ''
            
            if attempt >= retry_attempts:
                return stdout.decode('utf-8'), stderr.decode('utf-8')
            
            time.sleep(retry_delay)

    def get_http_download_token_by_request_id(self, access_token, request_id, filenames=None):
        """
        Get HTTP download token by request ID.

        Parameters:
        access_token (str): Access token for authorization.
        request_id (str): Request ID to identify the download request.
        filenames (list, optional): List of filenames to specify the files to download. Defaults to None.

        Returns:
        list: A list of dictionaries containing download tokens. Each dictionary contains:
        - "DownloadUrl" (str): The URL to download the file.
        - "FileName" (str): The name of the file to be downloaded.
        - "RequestId" (str): The unique identifier for the request.
        - "CompleteDate" (str): The date and time when the download was completed (YYYY/MM/DD HH:MM:SS).
        - "Token" (str): The token for authentication.

        Raises:
        ValueError: If the response JSON does not contain 'results' or if the request fails.
        """
        encoded_request_id = urllib.parse.quote(request_id)
        url = self.endpoint + f"/Customer/FEX/P2P/v1/requests/{encoded_request_id}/HTTPDownloadURLs"
        headers = {
            "Content-Type": "application/json; charset=utf-8",
            "Authorization": f"Bearer {access_token}"
        }
        if filenames:
            headers["filenames"] = "/".join(filenames)
        
        response = requests.get(url, headers=headers, verify=self.ssl_verify)
        if response.status_code == 200:
            response_data = response.json()
            if 'results' in response_data:
                return response_data['results']
            else:
                raise ValueError(f"Response JSON does not contain 'results': {response.text}")
        else:
            raise ValueError(f"Failed to get request details, status code: {response.status_code}, response: {response.text}")

    def download_files_by_http(self, access_token, request_id, filenames, target_dir, retry_delay=120, retry_attempts=3):
        """
        Download files via HTTP by request ID.

        Parameters:
        access_token (str): Access token for authorization.
        request_id (str): Request ID to identify the download request.
        filenames (list, optional): List of filenames to specify the files to download. Defaults to None.
        target_dir (str): Directory where the downloaded files will be saved.
        retry_delay (int, optional): Delay in seconds between retry attempts. Defaults to 120.
        retry_attempts (int, optional): Number of retry attempts if the download fails. Defaults to 3.
    
        Returns:
        list: A list of dictionaries containing information about files that failed to download.

        Raises:
        ValueError: If the download fails after the specified number of retry attempts.
        """
        # Check if target directory exists, if not, create it
        if not os.path.exists(target_dir):
            os.makedirs(target_dir)

        download_params = self.get_http_download_token_by_request_id(access_token, request_id, filenames)
        error_list = []
        
        for param in download_params:
            download_url = param["DownloadUrl"]
            file_name = param["FileName"]
            file_path = os.path.join(target_dir, file_name)
            token = param["Token"]
            download_file_headers = {
                "Authorization": f"Token {token}"
            }
            
            # If the file already exists, delete it to allow overwriting
            if os.path.exists(file_path):
                os.remove(file_path)

            for attempt in range(1, retry_attempts + 1):
                try:
                    response = requests.get(download_url, headers=download_file_headers, stream=True)
                    if response.status_code == 200:
                        with open(file_path, 'wb') as f:
                            for chunk in response.iter_content(chunk_size=1024*1024*2):
                                f.write(chunk)
                        break
                    else:
                        raise ValueError(f"Failed to download file {file_name}, status code: {response.status_code}")
                except Exception as e:
                    if attempt >= retry_attempts:
                        error_list.append({
                            "file_name": file_name,
                            "error": str(e)
                        })
                        break
                    time.sleep(retry_delay)
        
        return error_list


    