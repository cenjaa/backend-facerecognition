import numpy as np
from numpy.linalg import norm, svd


class R_PCA:
    def __init__(self, D, mu=None, lmbda=None):
        """
        Robust PCA (Principal Component Analysis) via Principal Component Pursuit (PCP).
        Decomposes matrix D into Low-rank matrix L and Sparse matrix S.
        """
        self.D = D
        self.S = np.zeros(self.D.shape)
        self.Y = np.zeros(self.D.shape)

        # Initialize mu (learning rate)
        if mu:
            self.mu = mu
        else:
            # Common heuristic: 1.25 / operator_norm(D) -> equivalent to norm(D, 2)
            self.mu = 1.25 / norm(self.D, 2)

        # Initialize lambda (sparsity weight)
        if lmbda:
            self.lmbda = lmbda
        else:
            # Standard heuristic: 1 / sqrt(max(dim1, dim2))
            self.lmbda = 1 / np.sqrt(np.max(self.D.shape))

    def fit(self, tol=1E-7, max_iter=1000):
        """
        Iterative solver using Augmented Lagrange Multiplier (ALM) method.
        """
        iter = 0
        err = np.inf
        Sk = self.S
        Yk = self.Y
        Lk = np.zeros(self.D.shape)

        print(f"   -> [RPCA] Starting decomposition on {self.D.shape} matrix...")

        while (err > tol) and (iter < max_iter):
            # 1. Update L (Low-Rank) using SVD Thresholding
            Lk = self.svd_thresholding(self.D - Sk + Yk/self.mu, 1/self.mu)

            # 2. Update S (Sparse) using Soft Thresholding
            Sk = self.soft_thresholding(self.D - Lk + Yk/self.mu, self.lmbda/self.mu)

            # 3. Update Y (Lagrange Multiplier) - The "Error" accumulator
            Z = self.D - Lk - Sk
            Yk = Yk + self.mu * Z

            # 4. Calculate Error (Frobenius Norm)
            err = norm(Z, 'fro') / norm(self.D, 'fro')

            iter += 1
            if iter % 50 == 0:
                print(f"   -> [RPCA] Iteration {iter}: Error {err:.7f}")

        self.L = Lk
        self.S = Sk
        print(f"   -> [RPCA] Converged after {iter} iterations. Final Error: {err:.7f}")
        return Lk, Sk

    def svd_thresholding(self, X, tau):
        """
        Shrinkage operator for singular values.
        """
        U, S, V = svd(X, full_matrices=False)
        return np.dot(U, np.dot(np.diag(self.soft_thresholding(S, tau)), V))

    def soft_thresholding(self, x, tau):
        """
        Standard soft thresholding operator.
        """
        return np.sign(x) * np.maximum(np.abs(x) - tau, 0)
