#!/bin/sh
# Enable Clique signer for node
geth attach /root/.ethereum/geth.ipc --exec "clique.propose('0x7df9a875a174b3bc565e6424a0050ebc1b2d1d82', true)"
