
import base64
import requests
import os
from src.temp_measure.utils.jwt import create_jwt

GCP_PROJECT_ID = 'temp-measure-dev'
PRIVATE_KEY_FILE_PATH = os.path.join(os.path.dirname(__file__), os.pardir, os.pardir, os.pardir, 'keys/rsa_private.pem')
ENCRYPTION_ALGO = 'RS256' 
CLOUD_REGION = 'europe-west1'
REGISTRY_ID = 'temp-sensors'
DEVICE_ID = 'rasp-pi-dht22'
HEADERS = {
    'authorization': 'Bearer {jwt_token}',
    'content-type': 'application/json',
    'cache-control': 'no-cache'
}
URL = f'https://cloudiotdevice.googleapis.com/v1/projects/{GCP_PROJECT_ID}/locations/{CLOUD_REGION}/registries/{REGISTRY_ID}/devices/{DEVICE_ID}:publishEvent'

class IOTInterface:
    def __init__(self):
        # TODO: Add token refreshing logic
        self._jwt = create_jwt(GCP_PROJECT_ID, PRIVATE_KEY_FILE_PATH, ENCRYPTION_ALGO).decode('ascii')  # decode again as need the string representation for the request

    def build_iot_post_request_payload(self, reading: str):
        encoded_temp = base64.b64encode(reading.encode('ascii')).decode('ascii') # decode again as need the string representation for the request
        return f'{{\"binary_data\": \"{encoded_temp}\"}}'


    def make_request(self, reading: str) -> None:
        HEADERS['authorization'] = HEADERS['authorization'].format(jwt_token=self._jwt)
        # TODO: add try/catch block
        payload = self.build_iot_post_request_payload(reading)

        resp = requests.post(URL, headers=HEADERS,data=payload)
        print(resp.text)
