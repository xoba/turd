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
                      rest) ;; the args
                    a))
                ((eq op 'quote) first)
                ((eq op 'cond) (evcon (cdr e) a))
                ((eq op 'list) (evlis (cdr e) a))
                ;; ((eq op 'display) (display (evlis (cdr e) a)))
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
                ;; "caaaaar" with 1 args (loaded)
                ((eq op 'caaaaar) (caaaaar (eval first a)))
                ;; "caaaadr" with 1 args (loaded)
                ((eq op 'caaaadr) (caaaadr (eval first a)))
                ;; "caaaar" with 1 args (loaded)
                ((eq op 'caaaar) (caaaar (eval first a)))
                ;; "caaadar" with 1 args (loaded)
                ((eq op 'caaadar) (caaadar (eval first a)))
                ;; "caaaddr" with 1 args (loaded)
                ((eq op 'caaaddr) (caaaddr (eval first a)))
                ;; "caaadr" with 1 args (loaded)
                ((eq op 'caaadr) (caaadr (eval first a)))
                ;; "caaar" with 1 args (loaded)
                ((eq op 'caaar) (caaar (eval first a)))
                ;; "caadaar" with 1 args (loaded)
                ((eq op 'caadaar) (caadaar (eval first a)))
                ;; "caadadr" with 1 args (loaded)
                ((eq op 'caadadr) (caadadr (eval first a)))
                ;; "caadar" with 1 args (loaded)
                ((eq op 'caadar) (caadar (eval first a)))
                ;; "caaddar" with 1 args (loaded)
                ((eq op 'caaddar) (caaddar (eval first a)))
                ;; "caadddr" with 1 args (loaded)
                ((eq op 'caadddr) (caadddr (eval first a)))
                ;; "caaddr" with 1 args (loaded)
                ((eq op 'caaddr) (caaddr (eval first a)))
                ;; "caadr" with 1 args (loaded)
                ((eq op 'caadr) (caadr (eval first a)))
                ;; "caar" with 1 args (loaded)
                ((eq op 'caar) (caar (eval first a)))
                ;; "cadaaar" with 1 args (loaded)
                ((eq op 'cadaaar) (cadaaar (eval first a)))
                ;; "cadaadr" with 1 args (loaded)
                ((eq op 'cadaadr) (cadaadr (eval first a)))
                ;; "cadaar" with 1 args (loaded)
                ((eq op 'cadaar) (cadaar (eval first a)))
                ;; "cadadar" with 1 args (loaded)
                ((eq op 'cadadar) (cadadar (eval first a)))
                ;; "cadaddr" with 1 args (loaded)
                ((eq op 'cadaddr) (cadaddr (eval first a)))
                ;; "cadadr" with 1 args (loaded)
                ((eq op 'cadadr) (cadadr (eval first a)))
                ;; "cadar" with 1 args (loaded)
                ((eq op 'cadar) (cadar (eval first a)))
                ;; "caddaar" with 1 args (loaded)
                ((eq op 'caddaar) (caddaar (eval first a)))
                ;; "caddadr" with 1 args (loaded)
                ((eq op 'caddadr) (caddadr (eval first a)))
                ;; "caddar" with 1 args (loaded)
                ((eq op 'caddar) (caddar (eval first a)))
                ;; "cadddar" with 1 args (loaded)
                ((eq op 'cadddar) (cadddar (eval first a)))
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
                ;; "cdaaaar" with 1 args (loaded)
                ((eq op 'cdaaaar) (cdaaaar (eval first a)))
                ;; "cdaaadr" with 1 args (loaded)
                ((eq op 'cdaaadr) (cdaaadr (eval first a)))
                ;; "cdaaar" with 1 args (loaded)
                ((eq op 'cdaaar) (cdaaar (eval first a)))
                ;; "cdaadar" with 1 args (loaded)
                ((eq op 'cdaadar) (cdaadar (eval first a)))
                ;; "cdaaddr" with 1 args (loaded)
                ((eq op 'cdaaddr) (cdaaddr (eval first a)))
                ;; "cdaadr" with 1 args (loaded)
                ((eq op 'cdaadr) (cdaadr (eval first a)))
                ;; "cdaar" with 1 args (loaded)
                ((eq op 'cdaar) (cdaar (eval first a)))
                ;; "cdadaar" with 1 args (loaded)
                ((eq op 'cdadaar) (cdadaar (eval first a)))
                ;; "cdadadr" with 1 args (loaded)
                ((eq op 'cdadadr) (cdadadr (eval first a)))
                ;; "cdadar" with 1 args (loaded)
                ((eq op 'cdadar) (cdadar (eval first a)))
                ;; "cdaddar" with 1 args (loaded)
                ((eq op 'cdaddar) (cdaddar (eval first a)))
                ;; "cdadddr" with 1 args (loaded)
                ((eq op 'cdadddr) (cdadddr (eval first a)))
                ;; "cdaddr" with 1 args (loaded)
                ((eq op 'cdaddr) (cdaddr (eval first a)))
                ;; "cdadr" with 1 args (loaded)
                ((eq op 'cdadr) (cdadr (eval first a)))
                ;; "cdar" with 1 args (loaded)
                ((eq op 'cdar) (cdar (eval first a)))
                ;; "cddaaar" with 1 args (loaded)
                ((eq op 'cddaaar) (cddaaar (eval first a)))
                ;; "cddaadr" with 1 args (loaded)
                ((eq op 'cddaadr) (cddaadr (eval first a)))
                ;; "cddaar" with 1 args (loaded)
                ((eq op 'cddaar) (cddaar (eval first a)))
                ;; "cddadar" with 1 args (loaded)
                ((eq op 'cddadar) (cddadar (eval first a)))
                ;; "cddaddr" with 1 args (loaded)
                ((eq op 'cddaddr) (cddaddr (eval first a)))
                ;; "cddadr" with 1 args (loaded)
                ((eq op 'cddadr) (cddadr (eval first a)))
                ;; "cddar" with 1 args (loaded)
                ((eq op 'cddar) (cddar (eval first a)))
                ;; "cdddaar" with 1 args (loaded)
                ((eq op 'cdddaar) (cdddaar (eval first a)))
                ;; "cdddadr" with 1 args (loaded)
                ((eq op 'cdddadr) (cdddadr (eval first a)))
                ;; "cdddar" with 1 args (loaded)
                ((eq op 'cdddar) (cdddar (eval first a)))
                ;; "cddddar" with 1 args (loaded)
                ((eq op 'cddddar) (cddddar (eval first a)))
                ;; "cdddddr" with 1 args (loaded)
                ((eq op 'cdddddr) (cdddddr (eval first a)))
                ;; "cddddr" with 1 args (loaded)
                ((eq op 'cddddr) (cddddr (eval first a)))
                ;; "cdddr" with 1 args (loaded)
                ((eq op 'cdddr) (cdddr (eval first a)))
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
            (car rest) ;; second
            (cadr rest))) ;; third
        (car e) ;; op
        (cadr e) ;; first
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
