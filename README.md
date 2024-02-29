# btc-service

Primarily built to help with Ordinox development and testing. The objective is as follows:

1. Help quickly inscribe brc20 token ops
2. Help transfer brc20 tokens between P2PKH addresses with minimal key management - as to mimic users

### How to use

1. Setup bitcoind, ord, bitcoin-cli, opi, brc-20 module and run them all
2. Create a "default.yaml" config file with the given sample. Enter all the values as requried
3. Build the program with `go build`
4. Create keypair(s) with the `./btc-service keypair` command

### Testing the transfer flow

1. Deploy a token
   `./btc-service brc20 inscribe-deploy [TOKEN] [SUPPLY] [ADDRESS]`
   This will deploy a token [TOKEN] with the supply of [SUPPLY]. The address given sends the inscription to that address.
   Note: This does not mean that the address has these tokens now. They still need to be minted.

2. Mint a token
   `./btc-service brc20 inscribe-mint [TOKEN] [AMT] [ADDRESS]`
   This will mint [AMT] number of [TOKEN] tokens into the given address.
   Note: Minting will only work if the token is deployed.
   Note: If you wish to transfer these tokens out of the given address, keep the privateKeyHex handy

3. Transfer a token
   This involves 2 steps.
   a. Inscribing a transfer inscription to the desired wallet address
   `./btc-service brc20 inscribe-transfer [FROM_ADDRESS] [TOKEN] [AMT]`
   This will inscribe a transfer inscription into the given wallet. This inscription is valid only if the given address actually has a balance of the [amt] number of tokens

   b. Sending the transfer inscription to the desired wallet address
   Grab the inscriptionid from the previous output, and then feed it into this command
   `./btc-service brc20 transfer [FROM_ADDRESS] [TO_ADDRESS] [TRANSFER_INSCRIPTION_ID] [AMT] [PRIVATE_KEY_HEX]`

4. Check balance at each step
   `btc-service brc20 balance [TOKEN] [ADDRESS]`

**IMPORTANT NOTE**
Make sure you generate some number of blocks before each step, or it will fail!
You can do that with
`./btc-service genblocks [AMT] [ADDRESS]`

`genblocks` command basically mines the blockchain and received block rewards (freely on regtest)
but these block rewards are vested for 100 blocks. So it's recommended that you first generate 100 blocks and then mine 1 at a time to save db space

It's also recommended that you genblocks into the address you want to transfer out Tokens/Inscriptions so that you have gas to pay towards the inscription transfer
