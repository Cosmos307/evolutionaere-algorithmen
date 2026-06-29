#!/usr/bin/env python3
"""
Plot convergence curves and boxplots from convergence CSV files.
Usage: python3 scripts/plot.py
Reads:  results/raw/convergence_*.csv
Writes: results/plots/convergence_<instance>.png
        results/plots/boxplot_<instance>.png
"""

import os
import re
import glob
import pandas as pd
import matplotlib.pyplot as plt

RAW_DIR = "results/raw"
PLOT_DIR = "results/plots"
FOCUS_MAX_EVAL = 10_000

os.makedirs(PLOT_DIR, exist_ok=True)

# --- load all convergence CSVs ---
pattern = os.path.join(RAW_DIR, "convergence_*.csv")
files = glob.glob(pattern)

if not files:
    print(f"No convergence CSVs found in {RAW_DIR}")
    exit(1)

records = []
known_algorithms = ["stochastic_hillclimber", "one_plus_one", "hillclimber", "ga"]
for path in files:
    fname = os.path.basename(path)
    # convergence_<instance>_<algorithm>_seed<N>.csv
    m = re.match(r"^convergence_(.+)_seed(\d+)\.csv$", fname)
    if not m:
        print(f"Skipping unrecognized filename: {fname}")
        continue
    body, seed = m.group(1), int(m.group(2))

    algorithm = None
    instance = None
    for alg in known_algorithms:
        suffix = "_" + alg
        if body.endswith(suffix):
            algorithm = alg
            instance = body[: -len(suffix)]
            break

    if algorithm is None or instance is None:
        print(f"Skipping filename without known algorithm suffix: {fname}")
        continue

    if algorithm == "one_plus_one":
        algorithm = "stochastic_hillclimber"

    df = pd.read_csv(path)
    df["instance"] = instance
    df["algorithm"] = algorithm
    df["seed"] = seed
    records.append(df)

all_data = pd.concat(records, ignore_index=True)

# --- plot per instance ---
for instance, inst_df in all_data.groupby("instance"):

    # 1. Convergence plot — avg best_cost per (algorithm, eval)
    fig, ax = plt.subplots(figsize=(10, 6))
    for alg, alg_df in inst_df.groupby("algorithm"):
        avg = alg_df.groupby("eval")["best_cost"].mean()
        ax.plot(avg.index, avg.values, label=alg)

    ax.set_title(f"Convergence — {instance}")
    ax.set_xlabel("Evaluations")
    ax.set_ylabel("Avg Best Cost")
    max_eval = int(inst_df["eval"].max())
    ax.set_xlim(0, min(FOCUS_MAX_EVAL, max_eval))
    ax.legend()
    ax.grid(True, alpha=0.3)
    plt.tight_layout()
    out = os.path.join(PLOT_DIR, f"convergence_{instance}.png")
    plt.savefig(out, dpi=150)
    plt.close()
    print(f"Saved {out}")

    # 2. Boxplot — final best_cost distribution per algorithm
    final = inst_df[inst_df["eval"] == inst_df["eval"].max()]
    # group final costs per algorithm
    algs = sorted(final["algorithm"].unique())
    data_per_alg = [final[final["algorithm"] == a]["best_cost"].values for a in algs]

    fig, ax = plt.subplots(figsize=(8, 6))
    ax.boxplot(data_per_alg, tick_labels=algs)
    ax.set_title(f"Final Cost Distribution — {instance}")
    ax.set_xlabel("Algorithm")
    ax.set_ylabel("Best Cost")
    ax.grid(True, alpha=0.3, axis="y")
    plt.tight_layout()
    out = os.path.join(PLOT_DIR, f"boxplot_{instance}.png")
    plt.savefig(out, dpi=150)
    plt.close()
    print(f"Saved {out}")
