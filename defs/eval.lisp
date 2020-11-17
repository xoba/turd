(defun eval (e a) 
  (cond
   ((atom e) (assoc e a))
   ((atom (car e))
    (cond
     ((eq (car e) 'quote)   (cadr e))
     ((eq (car e) 'atom)    (atom    (eval (cadr e) a)))
     ((eq (car e) 'eq)      (eq      (eval (cadr e) a)
			             (eval (caddr e) a)))
     ((eq (car e) 'display) (display (eval (cadr e) a))) 
     ((eq (car e) 'car)     (car     (eval (cadr e) a)))
     ((eq (car e) 'cdr)     (cdr     (eval (cadr e) a)))
     ((eq (car e) 'cons)    (cons    (eval (cadr e) a)
			             (eval (caddr e) a)))
     ((eq (car e) 'cond)    (evcon   (cdr e) a))
     ((eq (car e) 'list)    (evlis   (cdr e) a)) ; this should be a macro perhaps
     ('t (eval (cons (assoc (car e) a)
		     (cdr e))
	       a))))

   ((eq (caar e) 'macro) ; starting to play with a "macro" idea
    (cdr e))
   
   ((eq (caar e) 'label)
    (eval (cons (caddar e) (cdr e))
	  (cons (list (cadar e) (car e)) a)))
   ((eq (caar e) 'lambda)
    (eval (caddar e)
	  (append (pair (cadar e) (evlis (cdr e) a))
		  a)))))
