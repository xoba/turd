(defun eval (e a) 
  (cond
   ((atom e) (assoc e a))
   ((atom (car e))
    ((λ (op first rest) 
       ((λ (second third)
	  (cond
	   ((eq op 'funcall) (eval (cons
				    (eval first a) ;; the function
				    rest)          ;; the args
				   a))	   
	   ;; axioms:
	   ((eq op 'quote)   (cadr e))
	   ((eq op 'atom)    (atom    (eval first  a)))
	   ((eq op 'eq)      (eq      (eval first  a)
				      (eval second a)))
	   ((eq op 'car)     (car     (eval first  a)))
	   ((eq op 'cdr)     (cdr     (eval first  a)))
	   ((eq op 'cons)    (cons    (eval first  a)
				      (eval second a)))
	   ((eq op 'cond)    (evcon   (cdr e) a))
	   ((eq op 'list)    (evlis   (cdr e) a))

	   ;; time:
	   ((eq op 'after)   (after   (eval first a)
				      (eval second a)))
	   ;; arithmetic
	   ((eq op 'add)     (add     (eval first  a)
				      (eval second a)))
	   ((eq op 'inc)     (inc     (eval first  a)))
	   ((eq op 'sub)     (sub     (eval first  a)
				      (eval second a)))
	   ((eq op 'mul)     (mul     (eval first  a)
				      (eval second a)))
	   ((eq op 'exp)     (exp     (eval first  a)
				      (eval second a)
				      (eval third  a)))
	   ;; crypto:
	   ((eq op 'concat)  (concat  (eval first  a)
				      (eval second a)))
	   ((eq op 'hash)    (hash    (eval first  a)))
	   ((eq op 'hashed)  (hashed  (eval first  a)))
	   ((eq op 'newkey)  (newkey))
	   ((eq op 'pub)     (pub     (eval first   a)))
	   ((eq op 'sign)    (sign    (eval first   a)
				      (eval second  a)))
	   ((eq op 'verify)  (verify  (eval first   a)
				      (eval second  a)
				      (eval third   a)))
	   ;; debugging:
	   ((eq op 'display) (display (eval first a)))
	   ((eq op 'runes)   (runes (eval (cadr e) a)))
	   ((eq op 'err)     (err (eval (cadr e) a)))
	   
	   ((eq op 'test1)   (test1   (eval first a)))
	   ((eq op 'test2)   (test2   (eval first a)))
	   ((eq op 'test3)   (test3   (eval first a)))

	   ;; resolve an unknown op:
	   ('t (eval (cons (assoc op a)
			   (cdr e))
		     a))))
	(car  rest)   ;; second
	(cadr rest))) ;; third
     (car e)    ;; op
     (cadr e)   ;; first
     (cddr e))) ;; rest
   
   ;; initial macro concept, note the two evals:
   ((eq (caar e) 'macro)
    (eval (eval (cadddar e) (pair (caddar e) (cdr e))) a))
   
   ((eq (caar e) 'label)
    (eval (cons (caddar e) (cdr e))
	  (cons (list (cadar e) (car e)) a)))
   
   ((or
     (eq (caar e) 'lambda)
     (eq (caar e) 'λ))
    (cond
     ((atom (cadar e)) ; lexpr
      (eval (caddar e)
	    (cons (list (cadar e) (evlis (cdr e) a))
		  a)))
     ('t ; traditional lambda
      (eval (caddar e)
	    (append (pair (cadar e) (evlis (cdr e) a))
		    a)))))))
