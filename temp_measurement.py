import adafruit_dht
import time
import board
import base64
import requests
from utils.jwt import create_jwt

# --------- User Settings ---------
SENSOR_LOCATION_NAME = "Main Living Room"
SECONDS_BETWEEN_READS = 10
METRIC_UNITS = False
GCP_PROJECT_ID = 'temp-measure-dev'
PRIVATE_KEY_FILE_PATH = './keys/rsa_private.pem'
ENCRYPTION_ALGO = 'RS256' 
CLOUD_REGION = 'europe-west1'
REGISTRY_ID = 'temp-sensors'
DEVICE_ID = 'rasp-pi-dht22'
# ---------------------------------

dhtSensor = adafruit_dht.DHT22(board.D4)

# TODO: Add token refreshing logic
jwt_token = create_jwt(GCP_PROJECT_ID, PRIVATE_KEY_FILE_PATH, ENCRYPTION_ALGO).decode('ascii')  # decode again as need the string representation for the request

headers = {
    'authorization': f'Bearer {jwt_token}',
    'content-type': 'application/json',
    'cache-control': 'no-cache'
}
url = f'https://cloudiotdevice.googleapis.com/v1/projects/{GCP_PROJECT_ID}/locations/{CLOUD_REGION}/registries/{REGISTRY_ID}/devices/{DEVICE_ID}:publishEvent'


while True:
    try:
        humidity = dhtSensor.humidity
        temp_c = dhtSensor.temperature
        humidity = format(humidity,".2f")
    except RuntimeError as e:
        print(e)
        temp_c = "Error in measurement, skipping this beat"
        humidity = "Error in measurement, skipping this beat"
    print(SENSOR_LOCATION_NAME + " Temperature(C)", temp_c)
    print(SENSOR_LOCATION_NAME + " Humidity(%)", humidity)

    # TODO: Send payload as structured JSON so it's easier to parse on the receiving end
    encoded_temp = base64.b64encode(f"temp is {temp_c} C".encode("ascii")).decode('ascii') # decode again as need the string representation for the request
    payload = f'{{\"binary_data\": \"{encoded_temp}\"}}'
    resp = requests.post(url, headers=headers,data=payload)
    print(resp.text)
    time.sleep(SECONDS_BETWEEN_READS)