(defun eval (e a) 
  (cond
   ((atom e) (assoc e a))
   ((atom (car e))
    (cond

     ;; somehow cxr's don't work in interpreted mode:
     ;;
     ;; ((iscxr (car e)) (cxr (car e) (eval (cadr  e) a)))
     ;; ((eq (car e) 'iscxr) (iscxr (eval (cadr  e) a)))
     ;;
     ;; maybe, we need an axiom to convert atoms to runes?
     
     ;; axioms:
     ((eq (car e) 'quote)   (cadr e))
     ((eq (car e) 'atom)    (atom    (eval (cadr  e) a)))
     ((eq (car e) 'eq)      (eq      (eval (cadr  e) a)
			             (eval (caddr e) a)))
     ((eq (car e) 'car)     (car     (eval (cadr  e) a)))
     ((eq (car e) 'cdr)     (cdr     (eval (cadr  e) a)))
     ((eq (car e) 'cons)    (cons    (eval (cadr  e) a)
			             (eval (caddr e) a)))
     ((eq (car e) 'cond)    (evcon   (cdr e) a))

     ;; arithmetic:
     ;;
     ;; btw, we should be able to replace this section with
     ;; something like "(twoargs 'plus 'minus 'mult)",
     ;; which would be expanded into the following six lines:
     ;;
     ((eq (car e) 'plus)    (plus    (eval (cadr  e) a)
			             (eval (caddr e) a)))
     ((eq (car e) 'minus)   (minus   (eval (cadr  e) a)
			             (eval (caddr e) a)))
     ((eq (car e) 'mult)    (mult    (eval (cadr  e) a)
			             (eval (caddr e) a)))
     ;; time
     ((eq (car e) 'after)   (after   (eval (cadr  e) a)
				     (eval (caddr e) a)))
     ;; crypto
     ((eq (car e) 'concat)  (concat  (eval (cadr  e) a)
				     (eval (caddr e) a)))
     ((eq (car e) 'hash)    (hash    (eval (cadr  e) a)))
     ((eq (car e) 'newkey)  (newkey))
     ((eq (car e) 'pub)     (pub     (eval (cadr  e)  a)))
     ((eq (car e) 'sign)    (sign    (eval (cadr  e)  a)
				     (eval (caddr e)  a)))
     ((eq (car e) 'verify)  (verify  (eval (cadr  e)  a)
				     (eval (caddr e)  a)
				     (eval (cadddr e) a)))
     ;; debug
     ((eq (car e) 'display) (display (eval (cadr  e) a)))
     ((eq (car e) 'runes)   (runes (eval (cadr e) a)))
     ((eq (car e) 'err)     (err (eval (cadr e) a)))
     
     ((eq (car e) 'list)    (evlis   (cdr e) a))
     
     ;; unknown op
     ('t (eval (cons (assoc (car e) a)
		     (cdr e))
	       a))))
   
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
