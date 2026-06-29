# Projektdokumentation

## Überblick
Dieses Projekt vergleicht drei Suchverfahren für das **Set Covering Problem (SCP)**:

- Hillclimber
- Stochastic Hillclimber
- Genetischer Algorithmus

Ziel ist, eine möglichst günstige Auswahl von Spalten zu finden, sodass alle Zeilen abgedeckt sind.

## Was wird hier behandelt?
Eine Lösung ist gültig, wenn jede Zeile mindestens einmal abgedeckt ist. Bewertet werden die Lösungen nach:

- `uncovered`: wie viele Zeilen noch fehlen
- `cost`: Gesamtkosten der gewählten Spalten
- `redundancy`: unnötige, doppelte Abdeckung

Optional gibt es einen `Repair`-Schritt, der ungültige Lösungen wieder repariert und überflüssige Spalten entfernt.

## Aufbau
- `cmd/`: Startpunkt des Programms
- `internal/algorithm/`: Implementierung der drei Verfahren
- `internal/scp/`: Parser, Bewertung und Repair für SCP
- `internal/experiment/`: Konfiguration, Ausführung und Ergebnisformat
- `data/`: SCP-Instanzen
- `results/`: erzeugte CSV-Dateien und Plots
- `scripts/`: Python-Skript für die Auswertung

## Eingabedaten
Die Dateien in `data/` enthalten SCP-Instanzen. Eine Instanz beschreibt:

- wie viele Zeilen abgedeckt werden müssen
- wie viele Spalten zur Auswahl stehen
- welche Kosten jede Spalte hat
- welche Spalten welche Zeilen abdecken

## Projekt ausführen
Go-Version laut Projekt: `go 1.26.2`

Start:

```bash
go run ./cmd
```

Standardmäßig wird in `cmd/main.go` eine Instanz (`data/scp41.txt`) mit mehreren Seeds und mehreren Studien ausgeführt.

Die Ergebnisse landen in:

- `results/raw/<studie>/` für Konvergenz-CSV-Dateien
- Konsolenausgabe für die Lauf-Zusammenfassung

## Python-Skript ausführen
Das Skript `scripts/plot.py` erstellt Plots aus den Konvergenz-CSV-Dateien.

Falls nötig:

```bash
.venv/bin/pip install pandas matplotlib
```

Dann zum Beispiel so:

```bash
MPLCONFIGDIR=.matplotlib .venv/bin/python scripts/plot.py
```

Hinweis: Das Skript liest aktuell aus `results/raw`, die Experimente schreiben aber in Unterordner wie `results/raw/baseline`. Für Plots muss man also entweder `RAW_DIR` im Skript anpassen oder die gewünschten CSV-Dateien vorher nach `results/raw` legen.

## Verwendete Algorithmen
- **Hillclimber**: kippt pro Schritt genau ein Bit und nimmt bessere Nachbarn.
- **Stochastic Hillclimber**: mutiert mehrere Bits mit Wahrscheinlichkeit `p = PMutFactor / n`.
- **GA**: Population, Tournament Selection, Uniform Crossover, Mutation, Ersetzung des schlechtesten Individuums.

## Wichtige Parameter
Die wichtigsten Einstellungen stehen in `internal/experiment/config.go`:

Mögliche Config-Einstellungen sind:

- `InstancePaths []string`: Liste der Instanz-Dateien, z. B. `data/scp41.txt`
- `Seeds []int64`: Liste von Seeds für mehrere Wiederholungen
- `Budget int`: maximale Anzahl an Fitness-Auswertungen pro Lauf
- `UseRepair bool`: `true` oder `false`, je nachdem ob Repair verwendet werden soll
- `PInit float64`: Startwahrscheinlichkeit für eine gewählte Spalte, typisch zwischen `0.0` und `1.0`
- `PMutFactor float64`: Faktor für die Mutationsrate, daraus wird `PMutFactor / n`
- `AcceptEqual bool`: nur für den Hillclimber, ob gleich gute Nachbarn akzeptiert werden
- `PopSize int`: Populationsgröße des genetischen Algorithmus
- `TournamentSize int`: Größe des Turniers bei der Selektion
- `CrossoverProb float64`: Wahrscheinlichkeit für Crossover, typisch zwischen `0.0` und `1.0`
- `LogConvergence bool`: `true` oder `false`, ob Konvergenz mitgeschrieben wird
- `LogInterval int`: nach wie vielen Auswertungen ein CSV-Eintrag geschrieben wird
- `OutputDir string`: Zielordner für die erzeugten CSV-Dateien

Aktuell werden in `cmd/main.go` unter anderem diese Werte verwendet:

- `Budget = 50000`
- `UseRepair = true`
- `PInit = 0.1`
- `AcceptEqual = false`
- `PMutFactor = 1.0`
- `PopSize = 30`
- `TournamentSize = 3`
- `CrossoverProb = 0.7`
- `LogConvergence = true`
- `LogInterval = 100`

## Repair
Der `Repair`-Schritt macht eine Lösung wieder gültig, falls Zeilen nicht mehr abgedeckt sind. Danach werden überflüssige Spalten entfernt. So bleibt die Suche näher an gültigen Lösungen.

## Ergebnislauf
Ein Ergebnislauf ist immer:

`Instanz x Seed x Algorithmus`

Standardstudien in `cmd/main.go`:

- `baseline`
- `pmut_0.5_over_l`
- `pmut_2.0_over_l`
- `ga_pop_50`
- `ga_tournament_2`
- `ga_crossover_0.9`

Pro Lauf werden u. a. diese Werte ausgegeben:

- `best_cost`
- `uncovered`
- `redundancy`
- `runtime_ms`
- `gap_percent`

Zusätzlich werden Konvergenz-CSV-Dateien geschrieben. Sie enthalten pro Lauf Werte wie:

- `eval`
- `best_cost`
- `uncovered`
- `redundancy`
- `score`

Damit kann man die drei Verfahren und verschiedene Parameter direkt vergleichen. Die Plots zeigen den Verlauf der besten Kosten und die Verteilung der Endergebnisse.
