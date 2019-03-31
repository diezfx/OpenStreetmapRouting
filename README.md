# Backend


## Installation
Das Backend is in go geschrieben. Da Go modules benutzt wird, wird mindestens version 1.11 benötigt.
Es sollte alle Abhängigkeiten automatisch herunterladen.
Starten mit: go run main.go

In der res/config.yaml kann die verwendete Datei angegben werden und eingestellt werden obs neu berechnet wird oder ein vorhander graph geladen wird.


localhost:8000



Problemstellung:
    Finde den kürzesten Weg unter der Einschränkung, dass ein Auto nur eine bestimmte Reichweite hat. Wird eine Tankstelle angefahren kann die Reichweite erneuert werden. Es sollen jedoch so wenig wie möglich Tankstellen verwendet werden.
    -> Ziel:kürzeste Weg mit den wenigsten Tankstellen



## Einlesen:
1. Auslesen der Wege
2. Auslesen der Knoten
3. Vorbereiten zum Einsatz
   1. damit es effektiv im Speicher liegt, werden Ids von 0-n vergeben
   2. offset array berechnen
   3. Kosten berechnen
   4. Grid erstellen zum schnellen finden von Punkten
4. Tankstellen werden in ein extra feld gelesen
   1. Dem Haupgraph werden alle Tankstellenknoten hinzugefügt
   2. Kante und Kante zurück zum nächstliegenden Knoten im Haupgraph werden hinzugefügt
   3. Offsetliste neu berechnen & grid für tankstellen




## angepasster Dijkstra
1. CalcStationDijkstraSnapshots ist hauptfunktion
Starte Dijkstra 1:n und schaue ob Ziel erreichbar ist
2. Möglichkeiten:
      1. Falls erreichbar gebe route zurück
      1. Andernfalls suche die Tankstellen, die am Luftlinie am nächsten zum Ziel liegen
3. Speichere die Tankstellen mit dem bisherigen Weg und Kosten in einem Stack und schreibe alle auf eine "Blacklist"(damit nicht nocheinmal besucht werden)
4. Nehme element aus Stack und gehe zu Schritt 1
5. Falls zu viele Stationen(config) besucht wurden breche ab


## Andere Funktionen:
1. normaler Dijkstra
2. liste alle erreichbare Kanten in einem Bereich optional mit Reichweite
3. liste alle erreichbaren Tankstellen in einem Bereich optional mit Reichweite
4. Profiling mit pprof ist eingebunden


### Bemerkungen: 
1. normale Dijkstra bricht nicht ab, da für die anderen Zwecke nicht erforderlich
2. Dijstra basiert auf dem kürzesten Weg



### Alternative Betrachtungen:
1. Eine andere Möglichkeit wäre gewesen alle Tankstellen mit allen anderen Tankstellen zu verbinden
   Danach wäre die optimale Route findbar gewesen mit einem einzigen Dijkstra auf dem modifizierten Graph
   Aber: teure precomputation und viel Speicher


## Evaluation
Durch den gewählten greedy Ansatz kann das Resultat beliebig suboptimal sein. Im Vergleich zum normalen Dijkstra erhöht sich die Zeit abhängig von den möglichen Tankstellen sehr stark. 
Zusätzlich gibt es Probleme mit der Anbindung der Tankstellen an den Straßengraph bzw. viele nicht erreichbare Routen im Graph im allgemeinen.

-> Viel Verbesserungspotential :D


# Frontend

## Installation
npm install
npm run serve
mit Browser localhost:8080 öffnen
vue.js wird verwendet, weil ich es testen wollte.


## Features
### Modus: Route
1. Mit einem Klick auf die Map kann ein Start/Ziel festgelegt werden.
2. Wird 0 als Reichweite angeben wird ohne Reichweitenbeschränkung geuscht, sonst mit.
3. Zoomt man in diesem Modus näher hin werden Tankstellen und die im Graph verwendeten Routen angezeigt.

### Modus: Stationsreach
1. Zeigt ausgehend von dem Startpunkt alle erreichbaren Kanten und Tankstellen an.
2. Ist die Reichweite gesetzt wird diese berücksichtigt.

Warnung: Aus Debugging gründen ist keine Beschränkung der Entferung oder Verringerung der Kanten implementiert. Zu weit herausgezoomt nur bedingt



