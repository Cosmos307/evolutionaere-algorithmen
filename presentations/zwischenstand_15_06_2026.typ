#set page(
  width: 16cm,
  height: 9cm,
  margin: (x: 1.0cm, y: 0.7cm),
  footer: context [
    #set text(size: 9pt, fill: luma(90))
    #align(center)[
      HTWK Leipzig | Folie #counter(page).display("1")/#counter(page).final().at(0)
    ]
  ],
)

#set text(font: "Arial", size: 12pt)
#set par(justify: false)

#let slide(title, body) = [
  #text(size: 14pt, weight: "bold")[#title]
  #v(0.12cm)
  #body
  #pagebreak()
]

#slide(
  [Set-Covering mit evolutionären Algorithmen],
  [
    #text(size: 12.5pt, weight: "semibold")[Zwischenstand und erste Ergebnisse]
    #v(0.1cm)
    - Thema: Kostenminimierung unter harter Abdeckungsbedingung
    - Vergleich: Hillclimber, (1+1)-EA und genetischer Algorithmus (GA)
    - Ziel: fairer Methodenvergleich bei gleichem Evaluationsbudget
  ],
)

#slide(
  [1) Was ist das Problem?],
  [
    - Beim Set-Covering müssen alle Zeilen abgedeckt sein
    - Jede gewählte Spalte kostet Geld
    - Ziel: minimale Gesamtkosten bei vollständiger Abdeckung
    - Darstellung: Bitvektor (Spalte an/aus)
    - Gültig ist nur `Uncovered = 0`
  ],
)

#slide(
  [2) Wie gehen wir vor?],
  [
    - Parser liest OR-Library-Instanzen (z. B. `data/scp41.txt`)
    - Gleiche Fitness für alle Verfahren (Reihenfolge):
      weniger `Uncovered` -> kleinere `Cost` -> kleinere `Redundancy`
    - Constraint-Handling per `Repair`:
      - Greedy Spalten hinzufügen bis alles abgedeckt ist
      - Danach redundante Spalten entfernen
    - Runner führt alle Algorithmen über mehrere Seeds aus
  ],
)

#slide(
  [3) Welche Algorithmen und wie funktionieren sie?],
  [
    #set text(size: 10.5pt)
    #table(
      columns: (1.4fr, 2.4fr, 2.2fr),
      inset: 2pt,
      stroke: luma(210),
      [Algorithmus], [Prinzip], [Wichtige Parameter],
      [Hillclimber], [1-Bit-Nachbar, akzeptiert Verbesserung], [`PInit`, `AcceptEqual`, `Budget`],
      [(1+1)-EA], [Mehrfach-Bitflip mit `pMut`, akzeptiert nicht schlechter], [`PInit`, `PMut`, `Budget`],
      [GA], [Population + Turnier + Crossover + Mutation], [`PopSize`, `TournamentSize`, `CrossoverProb`, `PMut`],
    )
    #v(0.08cm)
    - Alle drei laufen mit identischer Fitness und identischem Budget
  ],
)

#slide(
  [4) Run-Setup und Ergebnisse],
  [
    #set text(size: 11pt)
    #grid(
      columns: (1.05fr, 1.45fr),
      gutter: 0.35cm,
      [
        #text(weight: "semibold")[Konfiguration]
        - Instanz: `scp41.txt`
        - Seeds: `1..5`
        - Budget: `50_000`
        - `UseRepair = true`
      ],
      [
        #text(weight: "semibold")[Finale Kosten (gemessener Run)]
        #v(0.08cm)
        #table(
          columns: (2.1fr, 1.2fr, 1.1fr, 1.1fr),
          inset: 2.5pt,
          stroke: luma(210),
          [Algorithmus], [Avg], [Best], [Worst],
          [GA], [437.8], [435], [439],
          [(1+1)-EA], [442.0], [433], [449],
          [Hillclimber], [444.0], [433], [457],
        )
      ],
    )
    #v(0.08cm)
    - Bei allen Verfahren: `Uncovered = 0`
  ],
)

#slide(
  [5) Convergence-Plot (0 bis 10k): Was zeigt er?],
  [
    #set text(size: 10.5pt)
    - Fokus auf die frühe Suchphase mit den meisten Verbesserungen
    - Niedriger ist besser; der GA liegt meist vorn
    #v(0.04cm)
    #align(center)[
      #image("../results/plots/convergence_scp41.txt.png", height: 4.4cm)
    ]
  ],
)

#slide(
  [6) Boxplot: Was zeigt die Verteilung?],
  [
    #set text(size: 10.5pt)
    - Streuung der finalen Kosten über alle Seeds
    - GA: beste mittlere Lage und relativ stabil
    - HC/(1+1)-EA: teils starke Einzelruns, aber mehr Schwankung
    #v(0.04cm)
    #align(center)[
      #image("../results/plots/boxplot_scp41.txt.png", height: 3.7cm)
    ]
  ],
)

#slide(
  [7) Wie wurden die Ergebnisse erreicht?],
  [
    - Einheitlicher Runner und gleiches Budget für fairen Vergleich
    - Seed-basierte Wiederholungen für robustere Aussagen
    - Repair hält Lösungen gültig und stabilisiert die Suche
    - Logging in `results/raw` macht Ergebnisse reproduzierbar
    - `scripts/plot.py` erzeugt die Vergleichsgrafiken
  ],
)

#slide(
  [8) Fazit und nächste Schritte],
  [
    #text(weight: "semibold")[Fazit]
    - Die End-to-End-Pipeline ist lauffähig
    - Auf `scp41` ist der GA aktuell im Mittel am besten

    #v(0.1cm)
    #text(weight: "semibold")[Nächste Schritte]
    - Weitere Instanzen (`scp42`, `scp43`, ...) für belastbarere Aussagen
    - Mehr Seeds testen (z. B. 30+) und Runs optional parallelisieren
    - Kleine Parameterstudie je Algorithmus
    - `BestKnown` pflegen und Gap sauber auswerten
    - Externe Config-Datei statt harter Werte in `cmd/main.go`
  ],
)
