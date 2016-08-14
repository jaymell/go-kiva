import requests
url = "http://api.kivaws.org/v1/loans/newest.json"
resp = requests.get(url).json()

for page in range(2,resp['paging']['total']):
	payload = {'page': page }
	print(requests.get(url, params=payload).json())


