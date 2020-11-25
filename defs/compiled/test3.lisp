(defun test3 (x)
  (car x))

;; TODO: label func's signature needs to be func(...Exp) Exp
;; in order to be used recursively
;; 
;;  ((label f
;;	  (lambda (first rest)
;;	    (cond
;;	     ((eq first '0) (list first rest))
;;	     ('t (f (minus first '1) rest))))
;;	  (car x) (cdr x))))

