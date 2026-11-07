
import matplotlib.pyplot as plt
import numpy as np
import os

# Create output directory
os.makedirs("graphs", exist_ok=True)

# Benchmark results
algorithms = ["Round Robin", "Weighted RR", "Least Connections"]

throughput = [43.99, 43.97, 43.97]

avg_latency = [1.61, 1.62, 1.60]
p95_latency = [2.19, 2.37, 2.57]
p99_latency = [2.67, 3.17, 4.69]

error_rate = [0, 0, 0]


# GRAPH 1: Throughput Comparison
plt.figure(figsize=(9, 5))

bars = plt.bar(algorithms, throughput)

plt.title("Throughput Comparison")
plt.xlabel("Scheduling Algorithm")
plt.ylabel("Throughput (requests/second)")
plt.ylim(0, 50)

for bar, value in zip(bars, throughput):
    plt.text(
        bar.get_x() + bar.get_width() / 2,
        value + 0.5,
        f"{value:.2f}",
        ha="center"
    )

plt.tight_layout()
plt.savefig("graphs/throughput_comparison.png", dpi=300)
plt.close()


# GRAPH 2: Latency Comparison
x = np.arange(len(algorithms))
width = 0.24

plt.figure(figsize=(9, 5))

plt.bar(x - width, avg_latency, width, label="Average Latency")
plt.bar(x, p95_latency, width, label="P95 Latency")
plt.bar(x + width, p99_latency, width, label="P99 Latency")

plt.title("Latency Comparison")
plt.xlabel("Scheduling Algorithm")
plt.ylabel("Latency (milliseconds)")
plt.xticks(x, algorithms)
plt.legend()

plt.tight_layout()
plt.savefig("graphs/latency_comparison.png", dpi=300)
plt.close()


# GRAPH 3: Error Rate Comparison
plt.figure(figsize=(9, 5))

bars = plt.bar(algorithms, error_rate)

plt.title("HTTP Error Rate Comparison")
plt.xlabel("Scheduling Algorithm")
plt.ylabel("Error Rate (%)")
plt.ylim(0, 1)

for bar, value in zip(bars, error_rate):
    plt.text(
        bar.get_x() + bar.get_width() / 2,
        0.05,
        f"{value:.2f}%",
        ha="center"
    )

plt.tight_layout()
plt.savefig("graphs/error_rate_comparison.png", dpi=300)
plt.close()

print("All three graphs generated successfully!")
print("Check the graphs/ directory.")
