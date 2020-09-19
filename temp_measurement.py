import adafruit_dht
import time
import board
from utils.jwt import create_jwt

# --------- User Settings ---------
SENSOR_LOCATION_NAME = "Main Living Room"
SECONDS_BETWEEN_READS = 10
METRIC_UNITS = False
PROJECT_ID = 'temp-measure-dev'
PRIVATE_KEY_FILE_PATH = './keys/rsa_private.pem'
ENCRYPTION_ALGO = 'RS256' 
# ---------------------------------

dhtSensor = adafruit_dht.DHT22(board.D4)

jwt_token = create_jwt(PROJECT_ID, PRIVATE_KEY_FILE_PATH, ENCRYPTION_ALGO)
print(jwt_token)

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
    time.sleep(SECONDS_BETWEEN_READS)