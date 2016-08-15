## TODO:
* Should this use a Client object like snapd and others, and requests are made through methods
	of that object?
* Immediate methods called by application, e.g., GetLoansById, should not worry about
	decoding json by themselves
  * Instead, have GetResponse do the decoding into an empty interface
