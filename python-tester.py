import requests

url = "http://api.kivaws.org/v1/loans/newest.json"
payload = {'page': 10, 'per_page': 75 }
resp = requests.get(url, params=payload).json()

url = "http://api.kivaws.org/v1/lenders/newest.json"
#payload = {'page': 10, 'per_page': 75 }
resp = requests.get(url).json()

url = "http://api.kivaws.org/v1/methods.json"
#payload = {'page': 10, 'per_page': 75 }
resp = requests.get(url).json()

url = "http://api.kivaws.org/v1/lenders/search.json"
payload = {'q': 'shoes'}
resp = requests.get(url).json()

for page in range(2,resp['paging']['total']):
	payload = {'page': page }
	print(requests.get(url, params=payload).json())


