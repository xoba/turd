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

	   ((eq op 'quote)   (cadr e))
	   ((eq op 'cond)    (evcon   (cdr e) a))
	   ((eq op 'list)    (evlis   (cdr e) a))

{{.compiled}}
	  
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
