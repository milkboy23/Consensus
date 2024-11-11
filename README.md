# Consensus
The instructions apply to opening the service in Visual Studio Code.
1. Clone the repository to your own machine.
2. In Visual Studio Code, open split terminal. Open as many terminals as you wish. The number of terminals is number of clients/nodes in the system.
3. Run all the terminal(s) by firstly running the go-file followed by the current id which is the port number starting from 0. The command should be formatted as such: "go run client.go -id 0" and so on. On the last client/terminal you must also write "-s true". This assigns the token to that client to begin with.

This prints constantly prints a number of messages, describing the token going from one client to the next depended on requests.

4. To close the program ctrl + c or shut down the terminals.
