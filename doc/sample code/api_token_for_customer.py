import requests
from requests.auth import HTTPBasicAuth
import os
import json
from datetime import datetime, timedelta

class AccessTokenForCustomerManager:
    """
    AccessTokenManager is responsible for managing access and refresh tokens.
    It handles fetching new tokens from the API server and storing them locally.

    Attributes:
        token_path (str): The file path where the token data is stored.
        ssl_verify (bool): Whether to verify SSL certificates for HTTPS requests.
    """
    
    def __init__(self,  endpoint, token_path, ssl_verify=True):
        self.token_path = token_path
        self.ssl_verify = ssl_verify
        self.endpoint = endpoint

    def fetch_refresh_token(self, username, password, client_id, secret):
        """
        Fetches a new token from the api server.

        Args:
            username (str): The username for authentication.
            password (str): The password for authentication.
            client_id (str): The client ID for authentication.
            secret (str): The secret for authentication.

        Returns:
            dict: The response data containing the refresh token.

        Raises:
            requests.HTTPError: If the POST request fails with a non-200 status code.
        """
        auth = HTTPBasicAuth(client_id, secret)
        headers = {
            'Content-Type': 'application/x-www-form-urlencoded'
        }
        body = {
            'grant_type': 'password',
            'username': username,
            'password': password
        }
        response = requests.post(self.endpoint + "/oauth/Customer/MOL/v1/token", auth=auth, headers=headers, data=body, verify=self.ssl_verify)
        if response.status_code == 200:
            return response.json()
        else:
            response.raise_for_status()

    def get_refresh_token(self, username, password, client_id, secret):
        """
        Retrieves or generates a refresh token for the given credentials.

        Args:
            username (str): The username for authentication.
            password (str): The password for authentication.
            client_id (str): The client ID for authentication.
            secret (str): The secret for authentication.

        Returns:
            dict: The response data containing the refresh token.

        Raises:
            requests.HTTPError: If the POST request fails with a non-200 status code.
        """
        if os.path.exists(self.token_path):
            with open(self.token_path, 'r') as file:
                data = json.load(file)
                if 'refresh_token_expires_at' in data:
                    refresh_token_expires_at = datetime.fromisoformat(data['refresh_token_expires_at'])
                    if refresh_token_expires_at > datetime.now():
                        return data

        new_data = self.fetch_refresh_token(username, password, client_id, secret)
        current_time = datetime.now()
        new_data['expire_at'] = (current_time + timedelta(seconds=new_data['expires_in']) - timedelta(seconds=5)).isoformat()
        new_data['refresh_token_expires_at'] = (current_time + timedelta(seconds=new_data['refresh_token_expires_in']) - timedelta(seconds=5)).isoformat()
        with open(self.token_path, 'w') as file:
            json.dump(new_data, file)
        return new_data

    def fetch_access_token(self, refresh_token, client_id, secret):
        """
        Fetches a new access token from the server.

        Args:
            url (str): The URL to make the POST request for generating the access token.
            refresh_token (str): The refresh token for authentication.
            client_id (str): The client ID for authentication.
            secret (str): The secret for authentication.

        Returns:
            dict: The response data containing the access token.

        Raises:
            requests.HTTPError: If the POST request fails with a non-200 status code.
        """
        auth = HTTPBasicAuth(client_id, secret)
        headers = {
            'Content-Type': 'application/x-www-form-urlencoded'
        }
        body = {
            'grant_type': 'refresh_token',
            'refresh_token': refresh_token
        }
        response = requests.post(self.endpoint + "/oauth/Customer/MOL/v1/token", auth=auth, headers=headers, data=body, verify=self.ssl_verify)
        if response.status_code == 200:
            return response.json()
        else:
            response.raise_for_status()

    def get_access_token(self, username, password, client_id, secret):
        """
        Retrieves or generates an access token using the access token.

        Args:
            url (str): The URL to make the POST request for generating the access token.
            client_id (str): The client ID for authentication.
            secret (str): The secret for authentication.

        Returns:
            dict: The response data containing the access token.

        Raises:
            requests.HTTPError: If the POST request fails with a non-200 status code.
        """
        if os.path.exists(self.token_path):
            with open(self.token_path, 'r') as file:
                data = json.load(file)
                if 'expire_at' in data:
                    expire_at = datetime.fromisoformat(data['expire_at'])
                    if expire_at > datetime.now():
                        return data

        refresh_token_data = self.get_refresh_token(username, password, client_id, secret)
        refresh_token = refresh_token_data['refresh_token']
        response_data = self.fetch_access_token(refresh_token, client_id, secret)
        current_time = datetime.now()
        response_data['expire_at'] = (current_time + timedelta(seconds=response_data['expires_in']) - timedelta(seconds=5)).isoformat()
        response_data['refresh_token_expires_at'] = (current_time + timedelta(seconds=response_data['refresh_token_expires_in']) - timedelta(seconds=5)).isoformat()
        with open(self.token_path, 'w') as file:
            json.dump(response_data, file)
        return response_data