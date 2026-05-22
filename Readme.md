# Návod na použití `vault-tool.exe`

## Popis
Pomocná aplikace při implementaci ne-databázových GINIS modulů (APG/DKS/AIB).
Umí zašifrovat heslo do `dpapi:` tvaru


## -encrypt
Zašifruje heslo do `dpapi:` pro GINIS-vault. 

```ps1
.\vault-tool.exe -encrypt
```

## -decrypt 
rozšifrovat heslo z trezoru (interaktivní způsob) - zeptá se uživatele na `dpapi:ZASIFROVANO`.

``` ps1
.\vault-tool.exe -decrypt
```

## -decrypt soubor.txt
Rozšifruje heslo z trezoru na windows 2022 (neinteraktivní způsob).
V souboru `soubor.txt` musí být zašifrované heslo v `dpapi:tvaru`

Některé windows terminály totiž nepřijmou více než 256 znaků z clipboardu.
```ps1
.\vault-tool.exe -decrypt soubor.txt
```
