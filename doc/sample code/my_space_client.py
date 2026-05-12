import subprocess
import requests
import os
import time
import urllib
import base64
import urllib.parse

def encode(value):
    """
    Encode the given value by first applying Base64 encoding and then URL query encoding.

    This method is primarily used to encode paths for safe transmission in URLs.

    Parameters:
    value (str): The value to be encoded.

    Returns:
    str: The encoded value.
    """
    encoded_val = base64.b64encode(value.encode()).decode()
    return urllib.parse.quote(encoded_val)

class MySpaceClient:
    """
    A client for interacting with the MySpace API.

    This class provides methods to perform various operations such as creating, updating, deleting spaces,
    uploading and downloading files, and retrieving information about spaces and files.

    Methods:
    - get_my_space_list: Retrieves the list of spaces for the current account.
    - get_sub_folders_files: Retrieves the list of files and subfolders within a specified folder.
    - get_file_info: Retrieves information about a specified file.
    - create_folder: Creates a new folder under the specified parent path.
    - rename_folder_file: Renames an existing folder or file under the specified path.
    - delete_folder_file: Deletes an existing folder or file under the specified path.
    - upload_files_via_http: Uploads multiple files to the specified parent folder using HTTP.
    - download_files_via_http: Downloads a file from the server to the specified target file path using HTTP.
    - upload_files_via_aspera: Uploads multiple files to the specified parent folder using Aspera.
    - download_files_via_aspera: Downloads multiple files from the server to the specified target directory using Aspera.
    """

    def __init__(self, endpoint, ssl_verify):
        self.endpoint = endpoint
        self.ssl_verify = ssl_verify
    
    def get_my_space_list(self, access_token):
        """
        Get the list of spaces for the current account.

        This method retrieves the list of spaces associated with the account identified by the provided access token.

        Parameters:
        access_token (str): The access token for authorization.

        Returns:
        list: A list of dictionaries, each representing a space with the following fields:
            - SpaceId (int): The unique identifier of the space.
            - SpaceName (str): The name of the space.
            - Creator (str): The name of the person who created the space.
            - Description (str): A description of the space.
            - Members (list): A list of members in the space, each represented by a dictionary with an Email field.
            - Owned (bool): Whether the current account owns the space.
            - LastModified (str): The timestamp of the last modification to the space.
            - Created (str): The timestamp of when the space was created.
            - MemberDelete (bool): Whether members can delete the space.
            - MemberModify (bool): Whether members can modify the space.
            - Owners (list): A list of owners of the space, each represented by a dictionary with an Email field.
            - RetentionDays (int): The number of days the space will be retained.
            - IpWhiteList (list or None): A list of IP addresses allowed to access the space, or None if not set.
            - Visitors (list or None): A list of visitors to the space, or None if not set.

        Example response:
        [
            {
                "SpaceId": 22222,
                "SpaceName": "xx",
                "Creator": "xxx",
                "Description": "xx",
                "Members": [
                {
                    "Email": "xx"
                }
                ],
                "Owned": true,
                "LastModified": "2024/09/26 11:03:32",
                "Created": "2024/09/26 11:03:32",
                "MemberDelete": false,
                "MemberModify": false,
                "Owners": [
                {
                    "Email": "xxx"
                }
                ],
                "RetentionDays": 180,
                "IpWhiteList": null,
                "Visitors": null
            }
        ]

        Raises:
        ValueError: If the space list could not be retrieved.
        """
        url = self.endpoint + "/Customer/FEX/MySpace/v1/spaces"
        headers = {
            "Content-Type": "application/json; charset=utf-8",
            "Authorization": f"Bearer {access_token}"
        }
        response = requests.get(url, headers=headers, verify=self.ssl_verify)
        resp_body_json = response.json()
        if 'results' not in resp_body_json:
            raise ValueError("Failed to get space list")
        return resp_body_json["results"];
    
    def get_sub_folders_files(self, access_token, path):
        """
        Get the files and subfolders within a specified folder.

        This method retrieves the list of files and subfolders under the folder identified by the provided path.

        Parameters:
        access_token (str): The access token for authorization.
        path (str): The path of the folder to retrieve contents from. The path will be Base64 encoded and then URL encoded.

        Example input:
        path = "/test2/1232"

        Returns:
        dict: A dictionary containing the count of items and a list of items (files and subfolders) with the following fields:
            - Name (str): The name of the file or subfolder.
            - Path (str): The full path of the file or subfolder.
            - Modified (str): The timestamp of the last modification.
            - Type (str): The type of the item ("File" or "Folder").
            - Size (int): The size of the file in bytes (only for files).
            - Count (int): The number of items within the folder (only for folders).

        Example output:
        {
        "Count": 4,
        "DataList": [
            {
            "Name": "xx",
            "Path": "xx",
            "Modified": "2023/07/31 17:39:48",
            "Type": "File",
            "Size": 4058079,
            "Count": 4
            }
        ]
        }

        Raises:
        ValueError: If the retrieval of subfolders and files fails.
        """
        url = self.endpoint + f"/Customer/FEX/MySpace/v1/foldersAndFiles?path={encode(path)}"
        headers = {
            "Content-Type": "application/json; charset=utf-8",
            "Authorization": f"Bearer {access_token}"
        }
        response = requests.get(url, headers=headers, verify=self.ssl_verify)
        resp_body_json = response.json()
        if 'results' not in resp_body_json:
            raise ValueError(f"Failed to get sub folder/file {response.text}")
        return resp_body_json["results"]
    
    def get_file_info(self, access_token, path):
        """
        Get information about a specified file.

        This method retrieves the details of a file identified by the provided path.

        Parameters:
        access_token (str): The access token for authorization.
        path (str): The path of the file to retrieve information from. The path will be Base64 encoded and then URL encoded.

        Example input:
        path = "/test2/1232/xxx"

        Returns:
        dict: A dictionary containing the details of the file with the following fields:
            - Name (str): The name of the file.
            - Path (str): The full path of the file.
            - Modified (str): The timestamp of the last modification.
            - Type (str): The type of the item ("File").
            - Size (int): The size of the file in bytes.
            - Count (int): The number of items within the folder (only for folders).

        Example output:
        {
        "Count": 4,
        "DataList": [
            {
            "Name": "xxx",
            "Path": "xxx",
            "Modified": "2023/07/31 17:39:48",
            "Type": "File",
            "Size": 4058079,
            "Count": 4
            }
        ]
        }

        Raises:
        ValueError: If the retrieval of file information fails.
        """
        result = self.get_sub_folders_files(access_token, path);
        if "DataList" not in result:
            raise ValueError(f"Failed to get sub folder/file {result}")
        return result["DataList"][0]
    
    def create_folder(self, access_token, parent_path, folder_name):
        """
        Create a folder under the specified parent path.

        This method creates a new folder using the provided parent path and folder name.

        Parameters:
        access_token (str): The access token for authorization.
        parent_path (str): The path of the parent folder where the new folder will be created. The path will be Base64 encoded and then URL encoded.
        folder_name (str): The name of the new folder to be created.

        Example input:
        parent_path = "/test2/1232"
        folder_name = "NewFolder"

        Returns:
        bool: True if the folder was created successfully, False otherwise.

        Raises:
        ValueError: If the folder creation fails.
        """
        url = self.endpoint + f"/Customer/FEX/MySpace/v1/folders?path={encode(parent_path)}"
        headers = {
            "Content-Type": "application/json; charset=utf-8",
            "Authorization": f"Bearer {access_token}"
        }
        data = {
            "name": folder_name
        }
        response = requests.post(url, headers=headers, json=data, verify=self.ssl_verify)
        resp_body_json = response.json()
        if 'results' not in resp_body_json or 'Success' not in resp_body_json["results"]:
            raise ValueError(f"Failed to create folder {response.text}")
        return resp_body_json["results"]["Success"]

    def rename_folder_file(self, access_token, current_path, new_name):
        """
        Rename a folder or file under the specified path.

        This method renames an existing folder or file identified by the provided current path.

        Parameters:
        access_token (str): The access token for authorization.
        current_path (str): The current path of the folder or file to be renamed. The path will be Base64 encoded and then URL encoded.
        new_name (str): The new name for the folder or file.

        Example input:
        current_path = "/test2/1232/OldName"
        new_name = "NewName"

        Returns:
        bool: True if the folder or file was renamed successfully, False otherwise.

        Raises:
        ValueError: If the renaming fails.
        """
        url = self.endpoint + f"/Customer/FEX/MySpace/v1/foldersAndFiles?path={encode(current_path)}"
        headers = {
            "Content-Type": "application/json; charset=utf-8",
            "Authorization": f"Bearer {access_token}"
        }
        data = {
            "name": new_name
        }
        response = requests.put(url, headers=headers, json=data, verify=self.ssl_verify)
        resp_body_json = response.json()
        if 'results' not in resp_body_json or 'Success' not in resp_body_json["results"]:
            raise ValueError(f"Failed to rename {response.text}")
        return resp_body_json["results"]["Success"]

    def delete_folder_file(self, access_token, current_path):
        """
        Delete a folder or file under the specified path.

        This method deletes an existing folder or file identified by the provided current path.

        Parameters:
        access_token (str): The access token for authorization.
        current_path (str): The path of the folder or file to be deleted. The path will be Base64 encoded and then URL encoded.

        Example input:
        current_path = "/test2/1232/FolderOrFileToDelete"

        Returns:
        bool: True if the folder or file was deleted successfully, False otherwise.

        Raises:
        ValueError: If the deletion fails.
        """
        url = self.endpoint + f"/Customer/FEX/MySpace/v1/foldersAndFiles?path={encode(current_path)}"
        headers = {
            "Content-Type": "application/json; charset=utf-8",
            "Authorization": f"Bearer {access_token}"
        }
        response = requests.delete(url, headers=headers, verify=self.ssl_verify)
        resp_body_json = response.json()
        if 'results' not in resp_body_json or 'Success' not in resp_body_json["results"]:
            raise ValueError(f"Failed to delete path {response.text}")
        return resp_body_json["results"]["Success"]

    def upload_files_via_http(self, access_token, parent_folder_path, filePaths, chunk_size=1024 * 1024 * 100, chunk_retry_attempts=10):
        """
        Upload files to a specified folder via HTTP.

        This method uploads multiple files to the specified parent folder using HTTP.

        Parameters:
        access_token (str): The access token for authorization.
        parent_folder_path (str): The path of the parent folder where the files will be uploaded.
        filePaths (list): A list of file paths to be uploaded.
        chunk_size (int, optional): The size of each chunk to be uploaded in bytes. Defaults to 100 MB.
        chunk_retry_attempts (int, optional): The number of retry attempts for each chunk upload. Defaults to 10.

        Example input:
        access_token = "your_access_token"
        parent_folder_path = "/test2/1232"
        filePaths = ["path/to/file1.txt", "path/to/file2.txt"]

        Returns:
        list: A list of dictionaries containing any errors that occurred during the upload process. Each dictionary has the following fields:
            - file (str): The path of the file that encountered an error.
            - error (str): The error message.

        Raises:
        ValueError: If the initial request to get upload URLs fails or if any errors occur during the file upload process.
        """
        url = self.endpoint + "/Customer/FEX/MySpace/v1/files/HTTPUploadURLs"
        headers = {
            "Content-Type": "application/json; charset=utf-8",
            "Authorization": f"Bearer {access_token}"
        }
        data = {
            "ServerPath": parent_folder_path,
            "ClientPaths": [
                {
                    "path": file,
                    "size": os.path.getsize(file)
                } for file in filePaths
            ]
        }
        response = requests.post(url, headers=headers, json=data, verify=self.ssl_verify)
        resp_body_json = response.json()
        if 'results' not in resp_body_json:
            raise ValueError(f"Failed to delete path {response.text}")
        token = resp_body_json["results"]["token"];
        url = resp_body_json["results"]["url"];
        requestId = resp_body_json["results"]["requestId"];
        error_list = []

        for file in filePaths:
            try:
                self.upload_file_in_chunks(url, token, requestId, file, chunk_size, chunk_retry_attempts)
            except Exception as e:
                error_list.append({"file": file, "error": str(e)})

        if error_list:
            raise ValueError(f"Errors occurred during file upload: {error_list}")

        return error_list

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
                
                for _ in range(1, retry_attempts + 1): 
                    response = requests.post(url, headers=headers, files=files, data=data)
                    if response.status_code != 200:
                        time.sleep(1)
                    break;                    
                else:
                    raise ValueError(f"Failed to upload chunk {chunk_index + 1} of file {file_name} after 10 attempts: {response.text}")
    
    def check_upload_completion(self, access_token_provider, file_size, server_file_path, wait_for_completion=False, check_upload_completion_interval=60, check_upload_completion_attempts=10):
        """
        Check the completion status of an upload.

        This method checks if a file has been completely uploaded by comparing the file size on the server with the expected file size.

        Parameters:
        access_token_provider (callable): A callable that returns the access token for authorization.
        file_size (int): The expected size of the file in bytes.
        server_file_path (str): The path of the file on the server.
        wait_for_completion (bool, optional): Whether to wait for the upload to complete if it is not already complete. Defaults to False.
        check_upload_completion_interval (int, optional): The interval in seconds between status checks. Defaults to 60 seconds.
        check_upload_completion_attempts (int, optional): The number of attempts to check the upload status. Defaults to 10 attempts.

        Example input:
        access_token_provider = lambda: "your_access_token"
        file_size = 1024
        server_file_path = "/path/to/server/file"
        wait_for_completion = True
        check_upload_completion_interval = 60
        check_upload_completion_attempts = 10

        Returns:
        bool: True if the upload is complete and the file size matches, False otherwise.

        Raises:
        ValueError: If the file size on the server does not match the expected file size after the specified number of attempts.
        """
        access_token = access_token_provider()
        file_info = self.get_file_info(access_token, server_file_path)
        if file_info["Size"] == file_size:
            return True
        if wait_for_completion:
            attempt = 1
            while file_size != file_info["size"] and attempt <= check_upload_completion_attempts:
                time.sleep(check_upload_completion_interval)
                access_token = access_token_provider()
                file_info = self.get_file_info(access_token, server_file_path)
                attempt += 1
        return False

    def download_files_via_http(self, access_token, target_file_path, server_file_path, retry_delay=120, retry_attempts=3):
        """
        Download a file from the server via HTTP.

        This method downloads a file from the server to the specified target file path on the local machine.

        Parameters:
        access_token (str): The access token for authorization.
        target_file_path (str): The local path where the downloaded file will be saved.
        server_file_path (str): The path of the file on the server to be downloaded.
        retry_delay (int, optional): The delay in seconds between retry attempts. Defaults to 120 seconds.
        retry_attempts (int, optional): The number of retry attempts if the download fails. Defaults to 3 attempts.

        Example input:
        access_token = "your_access_token"
        target_file_path = "/local/path/to/downloaded/file"
        server_file_path = "/server/path/to/file"
        retry_delay = 120
        retry_attempts = 3

        Returns:
        None

        Raises:
        ValueError: If the download request fails or if the file cannot be downloaded after the specified number of attempts.
        """
        url = self.endpoint + "/Customer/FEX/MySpace/v1/files/HTTPDownloadURLs"
        headers = {
            "Content-Type": "application/json; charset=utf-8",
            "Authorization": f"Bearer {access_token}"
        }
        data = {
            "Paths": [base64.b64encode(server_file_path.encode()).decode()]
        }
        response = requests.post(url, headers=headers, json=data, verify=self.ssl_verify)
        resp_body_json = response.json()
        if 'results' not in resp_body_json:
            raise ValueError(f"Failed to create upload request {response.text}")
        download_param = resp_body_json["results"];
        url = download_param["Url"];
        token = download_param["Token"];

        if os.path.exists(target_file_path):
            os.remove(target_file_path)
        
        download_file_headers = {
            "Authorization": f"Token {token}"
        }
        for attempt in range(1, retry_attempts + 1):
            try:
                response = requests.get(url, headers=download_file_headers, stream=True)
                if response.status_code == 200:
                    with open(target_file_path, 'wb') as f:
                        for chunk in response.iter_content(chunk_size=8192):
                            f.write(chunk)
                    return None
                else:
                    raise ValueError(f"Failed to download file {target_file_path}, status code: {response.status_code}")
            except Exception as e:
                if attempt >= retry_attempts:
                    raise ValueError(f"Failed to download file {target_file_path} after {retry_attempts} attempts: {str(e)}")
                time.sleep(retry_delay)

    def upload_files_via_aspera(self, access_token, aspera_install_dir, parent_folder_path, file_paths, retry_delay=120, retry_attempts=3):
        """
        Upload files to a specified folder via Aspera.

        This method uploads multiple files to the specified parent folder using Aspera.

        Parameters:
        access_token (str): The access token for authorization.
        aspera_install_dir (str): The directory where Aspera is installed.
        parent_folder_path (str): The path of the parent folder where the files will be uploaded.
        file_paths (list): A list of file paths to be uploaded.
        retry_delay (int, optional): The delay in seconds between retry attempts. Defaults to 120 seconds.
        retry_attempts (int, optional): The number of retry attempts if the upload fails. Defaults to 3 attempts.

        Example input:
        access_token = "your_access_token"
        aspera_install_dir = "/path/to/aspera"
        parent_folder_path = "/server/path/to/parent/folder"
        file_paths = ["/local/path/to/file1.txt", "/local/path/to/file2.txt"]
        retry_delay = 120
        retry_attempts = 3

        Returns:
        tuple: A tuple containing the standard output and standard error of the Aspera command.

        Raises:
        FileNotFoundError: If the required Aspera files are missing.
        ValueError: If the upload request fails or if the Aspera command is missing in the response.
        """
        # Check if required Aspera files exist
        ascp_path = os.path.join(aspera_install_dir, 'bin', 'ascp.exe')
        openssh_path = os.path.join(aspera_install_dir, 'etc', 'asperaweb_id_dsa.openssh')
        if not os.path.exists(ascp_path) or not os.path.exists(openssh_path):
            raise FileNotFoundError("Required Aspera files are missing")

        url = self.endpoint + "/Customer/FEX/MySpace/v1/files/AsperaUploadCommands"
        headers = {
            "Content-Type": "application/json; charset=utf-8",
            "Authorization": f"Bearer {access_token}"
        }
        data = {
            "ServerPath": parent_folder_path,
            "ClientPaths": file_paths
        }
        response = requests.post(url, headers=headers, json=data, verify=self.ssl_verify)
        resp_body_json = response.json()
        if 'results' not in resp_body_json:
            raise ValueError(f"Failed to create upload request {response.text}")
        aspera_upload_param = resp_body_json["results"];
        if "aspera_cmd" not in aspera_upload_param:
            raise ValueError(f"aspera_cmd is missing in the response {response.text}")
        ascp_temp = aspera_upload_param["aspera_cmd"]
        cmd = ascp_temp.replace("{{ascp.exe path}}", ascp_path).replace("{{openssh file path}}", openssh_path)

        # Retry logic
        for attempt in range(1, retry_attempts + 1):
            process = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
            stdout, stderr = process.communicate()

            if process.returncode == 0:
                return stdout.decode('utf-8'), ''
            
            if attempt >= retry_attempts:
                return stdout.decode('utf-8'), stderr.decode('utf-8')
            
            time.sleep(retry_delay)
        
    def download_files_via_aspera(self, access_token, aspera_install_dir, server_file_paths, target_dir, retry_delay=120, retry_attempts=3):
        """
        Download files from the server via Aspera.

        This method downloads multiple files from the server to the specified target directory on the local machine using Aspera.

        Parameters:
        access_token (str): The access token for authorization.
        aspera_install_dir (str): The directory where Aspera is installed.
        server_file_paths (list): A list of file paths on the server to be downloaded.
        target_dir (str): The local directory where the downloaded files will be saved.
        retry_delay (int, optional): The delay in seconds between retry attempts. Defaults to 120 seconds.
        retry_attempts (int, optional): The number of retry attempts if the download fails. Defaults to 3 attempts.

        Example input:
        access_token = "your_access_token"
        aspera_install_dir = "/path/to/aspera"
        server_file_paths = ["/server/path/to/file1.txt", "/server/path/to/file2.txt"]
        target_dir = "/local/path/to/downloaded/files"
        retry_delay = 120
        retry_attempts = 3

        Returns:
        tuple: A tuple containing the standard output and standard error of the Aspera command.

        Raises:
        FileNotFoundError: If the required Aspera files are missing.
        ValueError: If the download request fails or if the Aspera command is missing in the response.
        """
        # Check if required Aspera files exist
        ascp_path = os.path.join(aspera_install_dir, 'bin', 'ascp.exe')
        openssh_path = os.path.join(aspera_install_dir, 'etc', 'asperaweb_id_dsa.openssh')
        if not os.path.exists(ascp_path) or not os.path.exists(openssh_path):
            raise FileNotFoundError("Required Aspera files are missing")
        
        url = self.endpoint + "/Customer/FEX/MySpace/v1/files/AsperaDownloadCommands"
        headers = {
            "Content-Type": "application/json; charset=utf-8",
            "Authorization": f"Bearer {access_token}"
        }
        data = {
            "Paths": server_file_paths
        }
        response = requests.post(url, headers=headers, json=data, verify=self.ssl_verify)
        resp_body_json = response.json()
        if 'results' not in resp_body_json:
            raise ValueError(f"Failed to create upload request {response.text}")
        aspera_download_param = resp_body_json["results"];
        if "aspera_cmd" not in aspera_download_param:
            raise ValueError("ascpargs is missing in the response")
        aspera_cmd = aspera_download_param["aspera_cmd"]
        
        cmd = aspera_cmd.replace("{ascp.exe path}", ascp_path).replace("{openssh file path}", openssh_path).replace("{download dir}", f"\"{target_dir}\"")

        # Retry logic
        for attempt in range(1, retry_attempts + 1):
            process = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
            stdout, stderr = process.communicate()

            if process.returncode == 0:
                return stdout.decode('utf-8'), ''
            
            if attempt >= retry_attempts:
                return stdout.decode('utf-8'), stderr.decode('utf-8')
            
            time.sleep(retry_delay)






