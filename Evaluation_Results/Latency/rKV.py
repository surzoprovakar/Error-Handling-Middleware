import matplotlib
matplotlib.rcParams['pdf.fonttype'] = 42
matplotlib.rcParams['ps.fonttype'] = 42
import matplotlib.pyplot as plt
import numpy as np
from matplotlib.ticker import MultipleLocator

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
hue = [5642.5, 6071.166667, 6928.666667, 7562.5, 8649.333333]
mvr_undo = [4963.833333, 5714.333333, 6260.35, 7149, 7725.333333]
stateRoll = [5442.333333, 5913.5, 7164.5, 7698.166667, 9160.833333]

# Set up positions for points on X-axis
index = np.arange(len(undo_operations))

# Plotting the line chart
# plt.figure(figsize=(8, 6))

plt.plot(index, hue, label=r'\texttt{HUE}', marker='s', markersize=3, linestyle='--', linewidth=1, color='green')
plt.plot(index, mvr_undo, label=r'\emph{MVR-Undo}', marker='D', markersize=3, linestyle=':', linewidth=1, color='blue')
plt.plot(index, stateRoll, label=r'\emph{StateRoll}', marker='o', markersize=3, linestyle='-.', linewidth=1, color='darkmagenta')

# Adding labels and title
plt.xlabel(r'\# of undo operations')
plt.ylabel(r'Latency (ms)')
# plt.title('Deviation vs Number of Undo Operations')

plt.xticks(index, undo_operations)
plt.yticks()
plt.gca().yaxis.set_major_locator(MultipleLocator(1000))

plt.legend()

# Display the plot
plt.grid(True)
plt.tight_layout()

# plt.gcf().set_size_inches(8, 4)
plt.gcf().set_size_inches(3, 2)
# Saving the figure as a PDF with specified size
plt.savefig('rKV-Latency.pdf', format='pdf', dpi=300, bbox_inches='tight')

plt.show()
