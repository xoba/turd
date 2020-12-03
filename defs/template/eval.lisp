;; {{.comment}}
(defun {{.defun}} {{.args}}
  
  {{.nextlambda_prefix}}
  
  (cond
   ((atom e) ({{.assoc}} e a))
   ((atom (car e))
    ((λ (op first rest)
	((λ (second third)

	    (cond ;; TODO: can we compile funcall instead?
	     ((eq op 'funcall) ({{.eval}} (cons
					   ({{.eval}} first a) ;; the function
					   rest)          ;; the args
					  a))
	     
	     ((eq op 'quote)   first)
	     ((eq op 'cond)    ({{.evcon}}   (cdr e) a))
	     ((eq op 'list)    ({{.evlis}}   (cdr e) a))
	     
	     {{.compiled}}
	     
	     ;; resolve an unknown op:
	     ('t ({{.eval}} (cons ({{.assoc}} op a)
				  (cdr e))
			    a))))


	 (car  rest)   ;; second
	 (cadr rest))) ;; third
     (car e)    ;; op
     (cadr e)   ;; first
     (cddr e))) ;; rest
   
   ;; initial macro concept, note the two evals:
   ((eq (caar e) 'macro)
    ({{.eval}} ({{.eval}} (cadddar e) (pair (caddar e) (cdr e))) a))
   
   ((eq (caar e) 'label)
    ({{.eval}} (cons (caddar e) (cdr e))
	       (cons (list (cadar e) (car e)) a)))
   
   ((or ;; two different notations for lambda
     (eq (caar e) 'lambda)
     (eq (caar e) 'λ))
    (cond ;; two different forms for lambda
     ((atom (cadar e)) ; lexpr form (lambda x ...)
      ({{.eval}} (caddar e)
		 (cons (list (cadar e) ({{.evlis}} (cdr e) a))
		       a)))
     ('t ; traditional form (lambda (x...) ...)
      ({{.eval}} (caddar e)
		 (append (pair (cadar e) ({{.evlis}} (cdr e) a))
			 a))))))
  
  {{.nextlambda_suffix}})
