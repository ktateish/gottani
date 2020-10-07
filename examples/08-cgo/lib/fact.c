#include <math.h>

long long fact(long long n) {
	if (n <= 1) {
		return n;
	}
	return n * fact(n - 1);
}

double pi() { return M_PI; }
