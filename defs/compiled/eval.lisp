;; THIS FILE IS AUTOGENERATED, DO NOT EDIT!
(defun eval (e a)
  
  
  
  (cond
   ((atom e) (assoc e a))
   ((atom (car e))
    ((λ (op first rest)
	((λ (second third)

	    (cond ;; TODO: can we compile funcall instead?
	     ((eq op 'funcall) (eval (cons
					   (eval first a) ;; the function
					   rest)          ;; the args
					  a))
	     
	     ((eq op 'quote)   first)
	     ((eq op 'cond)    (evcon   (cdr e) a))
	     ((eq op 'list)    (evlis   (cdr e) a))
	     
	     
;; "add" with 2 args (manual)
((eq op 'add) (add (eval first a) (eval second a)))

;; "after" with 2 args (manual)
((eq op 'after) (after (eval first a) (eval second a)))

;; "and" with 2 args (loaded)
((eq op 'and) (and (eval first a) (eval second a)))

;; "append" with 2 args (loaded)
((eq op 'append) (append (eval first a) (eval second a)))

;; "assoc" with 2 args (loaded)
((eq op 'assoc) (assoc (eval first a) (eval second a)))

;; "atom" with 1 args (axiom)
((eq op 'atom) (atom (eval first a)))

;; "caadr" with 1 args (loaded)
((eq op 'caadr) (caadr (eval first a)))

;; "caar" with 1 args (loaded)
((eq op 'caar) (caar (eval first a)))

;; "cadar" with 1 args (loaded)
((eq op 'cadar) (cadar (eval first a)))

;; "caddar" with 1 args (loaded)
((eq op 'caddar) (caddar (eval first a)))

;; "cadddar" with 1 args (loaded)
((eq op 'cadddar) (cadddar (eval first a)))

;; "caddddar" with 1 args (loaded)
((eq op 'caddddar) (caddddar (eval first a)))

;; "caddddr" with 1 args (loaded)
((eq op 'caddddr) (caddddr (eval first a)))

;; "cadddr" with 1 args (loaded)
((eq op 'cadddr) (cadddr (eval first a)))

;; "caddr" with 1 args (loaded)
((eq op 'caddr) (caddr (eval first a)))

;; "cadr" with 1 args (loaded)
((eq op 'cadr) (cadr (eval first a)))

;; "car" with 1 args (axiom)
((eq op 'car) (car (eval first a)))

;; "cdar" with 1 args (loaded)
((eq op 'cdar) (cdar (eval first a)))

;; "cddar" with 1 args (loaded)
((eq op 'cddar) (cddar (eval first a)))

;; "cdddar" with 1 args (loaded)
((eq op 'cdddar) (cdddar (eval first a)))

;; "cddr" with 1 args (loaded)
((eq op 'cddr) (cddr (eval first a)))

;; "cdr" with 1 args (axiom)
((eq op 'cdr) (cdr (eval first a)))

;; "concat" with 2 args (manual)
((eq op 'concat) (concat (eval first a) (eval second a)))

;; "cons" with 2 args (axiom)
((eq op 'cons) (cons (eval first a) (eval second a)))

;; "display" with 1 args (manual)
((eq op 'display) (display (eval first a)))

;; "eq" with 2 args (axiom)
((eq op 'eq) (eq (eval first a) (eval second a)))

;; "err" with 1 args (manual)
((eq op 'err) (err (eval first a)))

;; "eval" with 2 args (loaded)
((eq op 'eval) (eval (eval first a) (eval second a)))

;; "evcon" with 2 args (loaded)
((eq op 'evcon) (evcon (eval first a) (eval second a)))

;; "evlis" with 2 args (loaded)
((eq op 'evlis) (evlis (eval first a) (eval second a)))

;; "exp" with 3 args (manual)
((eq op 'exp) (exp (eval first a) (eval second a) (eval third a)))

;; "factorial" with 1 args (loaded)
((eq op 'factorial) (factorial (eval first a)))

;; "hash" with 1 args (manual)
((eq op 'hash) (hash (eval first a)))

;; "hashed" with 1 args (manual)
((eq op 'hashed) (hashed (eval first a)))

;; "inc" with 1 args (loaded)
((eq op 'inc) (inc (eval first a)))

;; "length" with 1 args (loaded)
((eq op 'length) (length (eval first a)))

;; "mul" with 2 args (manual)
((eq op 'mul) (mul (eval first a) (eval second a)))

;; "newkey" with 0 args (manual)
((eq op 'newkey) (newkey))

;; "next" with 1 args (loaded)
((eq op 'next) (next (eval first a)))

;; "not" with 1 args (loaded)
((eq op 'not) (not (eval first a)))

;; "null" with 1 args (loaded)
((eq op 'null) (null (eval first a)))

;; "or" with 2 args (loaded)
((eq op 'or) (or (eval first a) (eval second a)))

;; "pair" with 2 args (loaded)
((eq op 'pair) (pair (eval first a) (eval second a)))

;; "pub" with 1 args (manual)
((eq op 'pub) (pub (eval first a)))

;; "runes" with 1 args (manual)
((eq op 'runes) (runes (eval first a)))

;; "sign" with 2 args (manual)
((eq op 'sign) (sign (eval first a) (eval second a)))

;; "sub" with 2 args (manual)
((eq op 'sub) (sub (eval first a) (eval second a)))

;; "tassoc" with 3 args (loaded)
((eq op 'tassoc) (tassoc (eval first a) (eval second a) (eval third a)))

;; "test1" with 1 args (loaded)
((eq op 'test1) (test1 (eval first a)))

;; "test2" with 1 args (loaded)
((eq op 'test2) (test2 (eval first a)))

;; "test3" with 1 args (loaded)
((eq op 'test3) (test3 (eval first a)))

;; "teval" with 3 args (loaded)
((eq op 'teval) (teval (eval first a) (eval second a) (eval third a)))

;; "tevcon" with 3 args (loaded)
((eq op 'tevcon) (tevcon (eval first a) (eval second a) (eval third a)))

;; "tevlis" with 3 args (loaded)
((eq op 'tevlis) (tevlis (eval first a) (eval second a) (eval third a)))

;; "verify" with 3 args (manual)
((eq op 'verify) (verify (eval first a) (eval second a) (eval third a)))


	     
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
   
   ((or ;; two different notations for lambda
     (eq (caar e) 'lambda)
     (eq (caar e) 'λ))
    (cond ;; two different forms for lambda
     ((atom (cadar e)) ; lexpr form (lambda x ...)
      (eval (caddar e)
		 (cons (list (cadar e) (evlis (cdr e) a))
		       a)))
     ('t ; traditional form (lambda (x...) ...)
      (eval (caddar e)
		 (append (pair (cadar e) (evlis (cdr e) a))
			 a))))))
  
  )
