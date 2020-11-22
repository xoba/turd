; this is a script that has time and signature constraints
((lambda
   (input block trans)
   ((lambda (pub t0 time thash sig)
      (cond
       ((and
	 (after time t0)
	 (verify pub thash sig))
	(hash (concat 'AA input)))
       ('t ())))
    '{{.pub}} ; pub
    '{{.t0}} ; t0
    (assoc 'time block) ; time
    (assoc 'hash trans) ; thash
    (assoc '{{.pubname}} (assoc 'arguments trans)))) ; sig
