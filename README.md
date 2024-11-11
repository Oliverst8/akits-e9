# akits-e9

How to start the program. Our final solution is in the 3-excercise-ring folder

- Go into the project folder.
- Open n-1 (where n is the total number of nodes you want in the program) terminals and as arguments write two port numbers in each, one for where this node is, and one for the next one. Fx. In the first one we might give it port 5000 and port 5001 as arguments, and then in the second one, we give it 5001 and 5002, and continue this for as many nodes as we want, except the last one.
- Create one last node, where you link it to the first one you created, and then add a random argument to the end. In our example, if we wanted 5 nodes, the last one might be on port 5004, and we would give it port 5004 (where this node is) and port 5000 (where we started), and then a third random argument (this could be 3).
- Write go in every terminal, to start all the nodes.

To end the program, just interrupt the program.
