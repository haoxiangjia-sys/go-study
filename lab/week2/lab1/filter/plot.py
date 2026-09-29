import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
import seaborn as sns

# Removes the irrelevant information from the results.csv file
contents = open("results_utf8.csv", "r", encoding="utf-8-sig").read().split('\n')
with open("parsed_results.csv", 'w') as file:
    file.write('name,time,range\n')
    for line in contents:
        if 'Filter' in line:
            file.write(line + '\n')

# Read in the saved CSV data.
benchmark_data = pd.read_csv('parsed_results.csv', header=0, names=['name', 'time', 'range'])

# Go stores benchmark results in nanoseconds. Convert all results to seconds.
# 如果 CSV 里是 1.049e+08，除以 1e9 得到秒
df["time"] = df["time"].astype(float) / 1e9

# Use the name of the benchmark to extract the number of worker threads used.
#  e.g. "Filter/16-8" used 16 worker threads (goroutines).
# Note how the benchmark name corresponds to the regular expression 'Filter/\d+_workers-\d+'.
# Also note how we place brackets around the value we want to extract.
benchmark_data['threads'] = benchmark_data['name'].str.extract(r'Filter/(\d+)_workers-\d+').apply(pd.to_numeric)
benchmark_data['cpu_cores'] = benchmark_data['name'].str.extract(r'Filter/\d+_workers-(\d+)').apply(pd.to_numeric)
print(benchmark_data)

# Plot a bar chart.
ax = sns.barplot(data=benchmark_data, x='threads', y='time', hue='threads', palette='tab10', legend=False)

# Set descriptive axis lables.
ax.set(xlabel='Worker threads used', ylabel='Time taken (s)')

# Display the full figure.
plt.show()