# Scenario 1:
# All local updates have been propagated
# So just undo the specified updates

LD_LIBRARY_PATH=. ../ast_server 3 localhost:5083 Addresses3.txt Scn1/Actions3.txt


# Scenario 2:
# Some local updates did not have a chance to be merged 
# with all replicas but are being undone; 
# make a record not to merge the updates being undone
# with the local replicas.

# LD_LIBRARY_PATH=. ../ast_server 3 localhost:5083 Addresses3.txt Scn2/Actions3.txt


# Scenario 3:
# Some rejected local updates may need to be revisited 
# to check if the performed undo would allow these 
# rejected updates to succeed now. 

# LD_LIBRARY_PATH=. ../ast_server 3 localhost:5083 Addresses3.txt Scn3/Actions3.txt


# Without Middleware
# ../server 3 localhost:5083 Addresses3.txt Actions3.txt