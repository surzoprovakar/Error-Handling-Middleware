import matplotlib
matplotlib.rcParams['pdf.fonttype'] = 42
matplotlib.rcParams['ps.fonttype'] = 42
import matplotlib.pyplot as plt
import numpy as np

plt.rcParams.update({
    'text.usetex': True,
    'text.latex.preamble': r'\usepackage[T1]{fontenc}',
    'font.size': 11,           
    'axes.labelsize': 11,
    'xtick.labelsize': 10,
    'ytick.labelsize': 10,
    'legend.fontsize': 9,
})

# Data 
undo_operations = [100, 250, 500, 750, 1000]

# Throughput (ops/sec)
hue_throughput = [152.37, 141.65, 127.69, 102.56, 79.91]
mvr_throughput = [169.05, 165.31, 156.65, 134.26, 114.48]
stateRoll_throughput = [139.61, 130.96, 102.26, 78.65, 61.26]

# Set up positions for points on X-axis
index = np.arange(len(undo_operations))

# Plotting the line chart
# plt.figure(figsize=(8, 6))

plt.plot(index, hue_throughput, label=r'\texttt{HUE}', marker='s', markersize=3, linestyle='--', linewidth=1, color='green')
plt.plot(index, mvr_throughput, label=r'\emph{MVR-Undo}', marker='D', markersize=3, linestyle=':', linewidth=1, color='blue')
plt.plot(index, stateRoll_throughput, label=r'\emph{StateRoll}', marker='o', markersize=3, linestyle='-.', linewidth=1, color='darkmagenta')

# Adding labels and title
plt.xlabel(r'\# of undo operations')
plt.ylabel(r'Throughput (ops/s)')
# plt.title('Deviation vs Number of Undo Operations')

plt.xticks(index, undo_operations)
plt.yticks(np.arange(0, 351, 50))
# plt.gca().yaxis.set_major_locator(MultipleLocator(20))

plt.legend()

# Display the plot
plt.grid(True)
plt.tight_layout()

# plt.gcf().set_size_inches(8, 4)
plt.gcf().set_size_inches(3, 2)
# Saving the figure as a PDF with specified size
plt.savefig('FeatureStorm-Replica9.pdf', format='pdf', dpi=300, bbox_inches='tight')

plt.show()
