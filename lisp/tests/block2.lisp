((lambda
   (input block trans)
   ((lambda (pub t0 time thash sig)
      (cond
       ((and
	 (after time t0)
	 (verify pub thash sig))
	(hash (concat 'AA input)))
       ('t ())))
    'MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE1uNFEakLlxd1P/+LoiZDAOfN6BQIiNTcEmKAoOGc8ARkiixUm/mU6ILcpPJCVSGpsh7pGms0Fa3ydtOjHPHJig
    '2020-11-22T11:41:03.395Z
    (assoc 'time block)
    (assoc 'hash trans)
    (assoc '17ETd6TdQu/oY15kJpLNSno1/npodenwImyR6zD58uc (assoc 'arguments trans))))
 
 'dJXQ2VKTozqr2XqxN0HZqOIYwwhBzhWd83o/54mrQ+4
 '((height 1000) (time 2020-11-22T11:41:03.397Z))
 '((type normal) (hash xkSAUPEHJVOu+qUNro6CKAf/2iSZGa8SclWsTliR7bs) (inputs ((0 ((quantity 3) (script (lambda (input block trans) ((lambda (sig) (cond ((and (after (assoc 'time block) '2020-11-22T11:41:03.395Z) (verify 'MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEuhn47WIVlZ/dcFYRvvOSvo2o7TaKw6o4c+umVQo6o0eBVLhheu+QtLUw8UlrTCPqo5Q4f8kJAvx6IJSVE3ZBWA (assoc 'hash trans) sig)) (hash (concat 'AA input))) ('t ()))) (assoc 'lQwP0JBVefaQTLTA8zkjagSnNGOtFTD/knhjS/we1os (assoc 'arguments trans))))))) (1 ((quantity 4) (script (lambda (input block trans) ((lambda (sig) (cond ((and (after (assoc 'time block) '2020-11-22T11:41:03.395Z) (verify 'MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE1uNFEakLlxd1P/+LoiZDAOfN6BQIiNTcEmKAoOGc8ARkiixUm/mU6ILcpPJCVSGpsh7pGms0Fa3ydtOjHPHJig (assoc 'hash trans) sig)) (hash (concat 'AA input))) ('t ()))) (assoc '17ETd6TdQu/oY15kJpLNSno1/npodenwImyR6zD58uc (assoc 'arguments trans))))))))) (outputs ((0 ((quantity 7) (address /+9JMIr3u8rksMUkg2UzpZVMtp07EFhivCT44ZnN5DM))))) (content ()) (arguments ((lQwP0JBVefaQTLTA8zkjagSnNGOtFTD/knhjS/we1os MEUCIQDuVq9Jh8uu7EtJTEjhvclhJfvWURk7A+YR+3TFPCnEfAIgBAd12i42+jK71+w8QKMQ8tc3fwAq2CiLSfiSJAycbVY) (17ETd6TdQu/oY15kJpLNSno1/npodenwImyR6zD58uc MEUCIBN0M6xHiYd6kk//OccLqgk90KORc8U7JLbdWT0MOMRdAiEAxSJbxSm8sakAYnM+w4FOZ0H1nE+PboaN6vDQxPv3SIc)))))

