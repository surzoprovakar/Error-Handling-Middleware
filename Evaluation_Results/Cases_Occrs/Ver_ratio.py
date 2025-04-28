import matplotlib
matplotlib.rcParams['pdf.fonttype'] = 42
matplotlib.rcParams['ps.fonttype'] = 42
import matplotlib.pyplot as plt

plt.rcParams.update({
    'text.usetex': True,
    'text.latex.preamble': r'\usepackage[T1]{fontenc}',
    'font.size': 9,           
    'axes.labelsize': 9,
    'xtick.labelsize': 9,
    'ytick.labelsize': 9,
    'legend.fontsize': 8,
})

# Data
cases = ['Case 1: w/ Recourse', 'Case 1: w/o Recourse', 'Case 2', 'Case 3', 'Case 4']
raw_values = [27.06, 13.94, 19, 26, 14]

# Normalize to 100%
total = sum(raw_values)
values = [(v / total) * 100 for v in raw_values]

# Colors and hatch patterns
edgecolors = ['steelblue', 'steelblue', 'indigo', 'darkolivegreen', 'indianred']
hatch_patterns = ['...', 'xx', '||', '**', '//']

# Figure setup
fig, ax = plt.subplots(figsize=(0.5, 2))  # Narrow & tall for 2-column layout

# Plot
bottom = 0
bar_width = 0.05  # Reduced width
for i in range(len(values)):
    ax.bar(0, values[i], bottom=bottom, edgecolor=edgecolors[i], fill=False,
           hatch=hatch_patterns[i], width=bar_width, label=cases[i])
    bottom += values[i]

# Axes
ax.set_xticks([0])
ax.set_xticklabels([''])
ax.set_ylim(0, 100)
ax.set_yticks([0, 25, 50, 75, 100])
ax.set_yticklabels([r'0\%', r'25\%', r'50\%', r'75\%', r'100\%'])

# Legend matching bar height
ax.legend(loc='center left', bbox_to_anchor=(1.2, 0.5), frameon=True,
          borderaxespad=0.0, handlelength=1.5, edgecolor='dimgray')

# Save
plt.savefig('case-distribution.pdf', format='pdf', bbox_inches='tight')