Redshirt2Crypt
==============
Decrypter and Encrypter for redshirt2 files, used by Introversion's
[Darwinia](https://store.steampowered.com/app/1500/Darwinia/).  Reads stdin and
decrypts/encrypts to stdout, based on the presence or absense of the
`redshirt2` magic header.

Based heavily on [this code](https://forums.introversion.co.uk/viewtopic.php?p=479470&sid=08b726f0e5eb1793912baa2e74a90ace#p479470).

Quickstart
----------
```bash
# Install the program, you'll need the Go compiler
go install github.com/magisterquis/redshirt2crypt

# Decrypt the user data file
redshirt2crypt < "$HOME/Library/Application Support/Darwinia/full-steam/users/YOU/game.txt" >$HOME/Desktop/tmp

# Make whatever changes
vi $HOME/Desktop/tmp

# Re-encrypt the user data file
redshirt2crypt > "$HOME/Library/Application Support/Darwinia/full-steam/users/YOU/game.txt" <$HOME/Desktop/tmp
```

