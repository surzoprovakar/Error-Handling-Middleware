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
hue = [245.4, 266.5, 280.9, 308.4, 385.5]
mvr_undo = [202.6, 223.8, 228.9, 259.5, 309.2]
stateRoll = [233.3, 255.6, 292.03, 316.6, 397.5]

# Set up positions for points on X-axis
index = np.arange(len(undo_operations))

# Plotting the line chart
# plt.figure(figsize=(8, 6))

plt.plot(index, hue, label=r'\texttt{HUE}', marker='s', markersize=3, linestyle='--', linewidth=1, color='green')
plt.plot(index, mvr_undo, label=r'\emph{MVR-Undo}', marker='D', markersize=3, linestyle=':', linewidth=1, color='blue')
plt.plot(index, stateRoll, label=r'\emph{StateRoll}', marker='o', markersize=3, linestyle='-.', linewidth=1, color='darkmagenta')

# Adding labels and title
plt.xlabel(r'\# of undo operations')
plt.ylabel(r'Peak Memory (MB)')
# plt.title('Deviation vs Number of Undo Operations')

plt.xticks(index, undo_operations)
plt.yticks()
# plt.yticks(np.arange(0, 501, 100))

plt.legend()

# Display the plot
plt.grid(True)
plt.tight_layout()

# plt.gcf().set_size_inches(8, 4)
plt.gcf().set_size_inches(3, 2)
# Saving the figure as a PDF with specified size
plt.savefig('rKV-Memory.pdf', format='pdf', dpi=300, bbox_inches='tight')

plt.show()