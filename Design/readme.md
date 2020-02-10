# SANNTIDSPROGRAMMERING

## Prioritet

>Hver enkelt node har et nummer 1,2,...,n. Master er nr 1, mens Backup er nr 2. Resterende heiser har nr 3 ,...,n.

## Kommunikasjon
>Alle heisene snakker sammen over TCP, med master som koordinerende aktør. All informasjon samles inn til master, sendes til backup og sendes deretter ut til alle noder. Ved en handling informeres master og backup FØR noen handling utføres.

## Vaktbikkje
> Fungerende master purrer gjevnlig på alle nodene. Dersom bekreftelse ikke mottas fra en node registreres denne som inaktiv. Handlingskøen fra denne noden vil da flyttes fra aktiv i den aktuelle noden til inaktiv, slik at andre noder kan ta over arbeidet. Dersom noden fortsetter arbeidet vil de aktive bestillinger kunne betjenes 2 ganger, men sett fra nettverket kun en gang.

## Toleranse
> Master sitter på enhver tid på all informasjon, og oppdaterer backup. Dersom master svikter vil backup ved hjelp av en timer registrere at den ikke får livstegn fra master. Da vil backup ta over som master, og gjevnlig forsøke å vekke master tilbake til live. Når tidligere master er på nett igjen fungerer den som backup. Alle heiser har samme programvare, slik at master kan gå fra 1 til 2 til 3 og nedover.
Sett fra nodenes perspektiv vil den eneste endringen være adressen ved et master-bytte.
> På denne måten vil en svikt hvor som helst i systemet ikke kunne ta ned hele driften fordi to noder til enhver tid sitter på en oversikt som kan gjenopptas av en backup ved eventuell svikt.. 