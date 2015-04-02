# averno
Small, local NNTP server using DHT as a peering mechanism.


I will  write 4 components:

first, the dht component:

>-dht : using the dht library, the server will be able to find other peers.

then the second component, NNTP

1. NNTP client: once another peer is found, using NNTP it will upload and download new contents from the peer.
2. NNTP server: interface to the user's client and to the peers, in order to exchange messages and groups.

the third module, UpNP
  
>It will try to open port 11119 on the router (nntp) and some UDP port for dht.


Basically i plan to have 5 threads :

1. Upnp opening the router and keeping it open
2. dht looking for peers
3. nntp listening 
4. one infinite loop checking lifespans and deleting old articles from filesystem.
5. one infinite loop spreading/downloading new messages and groups to/from all known peers, using NNTP client.

all files will be stored in local ~/news folder. Since of port 11119 , no **root** privileges SHALL be needed.

I will make it assuming Linux/Unix. If it doesn't works in windows, please restart your PC.

I will choose UPnP , dht and NNTP libraries (if any) from existing ones. I hope.

1. DHT: [http://godoc.org/github.com/anacrolix/torrent/dht](http://godoc.org/github.com/anacrolix/torrent/dht) or [https://github.com/nictuku/dht](https://github.com/nictuku/dht)
2. NNTP: [https://github.com/dustin/go-nntp](https://github.com/dustin/go-nntp)
3. UPnP: [https://github.com/jackpal/go-nat-pmp](https://github.com/jackpal/go-nat-pmp)



