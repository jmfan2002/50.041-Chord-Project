import hashlib
import numpy as np
import matplotlib.pyplot as plt

def hash_and_plot():
    # Generate 10000 random integers
    integers = np.random.randint(1, 10001, size=10000000)

    # Hash the integers using SHA256 and store the hashes
    hashes = [int(hashlib.sha256(str(i).encode('utf-8')).hexdigest()[:16], 16) for i in integers]

    # Plot the distribution of the hashes
    plt.hist(hashes, bins=100)
    plt.title("10 Million SHA256 hashes of random numbers")
    plt.show()

hash_and_plot()