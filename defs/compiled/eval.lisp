(defun eval (e a) 
  (cond
   ((atom e) (assoc e a))
   ((atom (car e))
    ((lambda (op first second third)
       (cond

	((eq op 'test1)   (test1   (eval first a)))
	((eq op 'test2)   (test2   (eval first a)))
	
	;; axioms:
	((eq op 'quote)   (cadr e))
	((eq op 'atom)    (atom    (eval first a)))
	((eq op 'eq)      (eq      (eval first a)
				   (eval second a)))
	((eq op 'car)     (car     (eval first a)))
	((eq op 'cdr)     (cdr     (eval first a)))
	((eq op 'cons)    (cons    (eval first a)
				   (eval second a)))
	((eq op 'cond)    (evcon   (cdr e) a))

	((eq op 'plus)    (plus    (eval first a)
				   (eval second a)))
	((eq op 'inc)     (plus    (eval first a) '1))
	((eq op 'minus)   (minus   (eval first a)
				   (eval second a)))
	((eq op 'mult)    (mult    (eval first a)
				   (eval second a)))
	((eq op 'exp)     (exp     (eval first  a)
				   (eval second  a)
				   (eval third a)))
	;; time:
	((eq op 'after)   (after   (eval first a)
				   (eval second a)))
	;; crypto:
	((eq op 'concat)  (concat  (eval first a)
				   (eval second a)))
	((eq op 'hash)    (hash    (eval first a)))
	((eq op 'newkey)  (newkey))
	((eq op 'pub)     (pub     (eval first  a)))
	((eq op 'sign)    (sign    (eval first  a)
				   (eval second  a)))
	((eq op 'verify)  (verify  (eval first  a)
				   (eval second  a)
				   (eval third a)))
	;; debug:
	((eq op 'display) (display (eval first a)))
	((eq op 'runes)   (runes (eval (cadr e) a)))
	((eq op 'err)     (err (eval (cadr e) a)))
	
	((eq op 'list)    (evlis   (cdr e) a))

	;; resolve an unknown op:
	('t (eval (cons (assoc op a)
			(cdr e))
		  a))))
     (car e)    ;; op
     (cadr e)   ;; first
     (caddr e)  ;; second
     (cadddr e) ;; third
     )) 
   
   ;; initial macro concept, note the two evals:
   ((eq (caar e) 'macro)
    (eval (display (eval (cadddar e) (pair (caddar e) (cdr e)))) a))
   
   ((eq (caar e) 'label)
    (eval (cons (caddar e) (cdr e))
	  (cons (list (cadar e) (car e)) a)))
   
   ((eq (caar e) 'lambda)
    (cond
     ((atom (cadar e)) ; lexpr
      (eval (caddar e)
	    (cons (list (cadar e) (evlis (cdr e) a))
		  a)))
     ('t ; traditional lambda
      (eval (caddar e)
	    (append (pair (cadar e) (evlis (cdr e) a))
		    a)))))))
