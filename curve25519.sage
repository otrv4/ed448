# create curve
ec = EllipticCurve(GF(2**255-19), [0,486662,0,1,0])
k.<a> = GF(2**255-19)
base_point = ec.lift_x(9)
print(base_point)
point_at_infinity = ec(0)

# all elements of order 4
G4s = ec.lift_x(1, True)
G4 = ec.lift_x(1) # just the first element

# the element of order 2
G2 = G4 + G4

P = ec.random_point()
print(P)
Q = ec.random_point()
print(Q)

Z = P + Q
print(Z)
