import adafruit_dht
import time
import board
import base64
import requests
from src.utils.jwt import create_jwt
from src.gcp.iot.request_builder import IOTInterface

# --------- User Settings ---------
SECONDS_BETWEEN_READS = 10
# ---------------------------------

def get_single_reading(dhtSensor):
    try:
        humidity = dhtSensor.humidity
        temp_c = dhtSensor.temperature
        humidity = format(humidity,".2f")
    except RuntimeError as e:
        print(e)
        temp_c = "Error in measurement, skipping this beat"
        humidity = "Error in measurement, skipping this beat"
    print("Temperature(C)", temp_c)
    print("Humidity(%)", humidity)

    return temp_c, humidity


def measure_temperature():
    # Initialise sensor interface and GCP IOT interface
    dhtSensor = adafruit_dht.DHT22(board.D4)
    iot_interface = IOTInterface()

    while True:
        temp_c, humidity = get_single_reading(dhtSensor)

        iot_interface.make_request(temp_c, humidity)
        
        time.sleep(SECONDS_BETWEEN_READS)