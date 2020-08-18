import itertools
import time

from utils.aima_utils import (
    Expr, expr
)

from utils.aima_logic import (
    FolKB, pl_resolution, standardize_variables, parse_definite_clause,
    unify, subst, is_var_symbol, constant_symbols, variables
)


class kb_engine(FolKB):
    def tell(self, expression):
        if isinstance(expression, list) or isinstance(expression, tuple):
            for e in expressions:
                self.tell(e)
        elif isinstance(expression, str):
            expression = expr(expression)
        assert isinstance(expression, Expr)
        super().tell(expression)

    def ask(self, statement):
        if isinstance(statement, str):
            statement = expr(statement)
        assert isinstance(statement, Expr)
        return super().ask(statement)

    def clear(self):
        self.clauses = []

    """
        Check if FolKB entails the query, alpha
    """
    def fol_fc_ask(self, alpha):
        if isinstance(alpha, str):
            alpha = expr(alpha)
        """A simple forward-chaining algorithm. [Figure 9.3]"""
        # TODO: Improve efficiency
        kb_consts = list({c for clause in self.clauses for c in constant_symbols(clause)})
        def enum_subst(p):
            query_vars = list({v for clause in p for v in variables(clause)})
            for assignment_list in itertools.product(kb_consts, repeat=len(query_vars)):
                theta = {x: y for x, y in zip(query_vars, assignment_list)}
                yield theta

        # check if we can answer without new inferences
        for q in self.clauses:
            phi = unify(q, alpha, {})
            if phi is not None:
                yield phi

        while True:
            new = []
            for rule in self.clauses:
                p, q = parse_definite_clause(rule)
                for theta in enum_subst(p):
                    if set(subst(theta, p)).issubset(set(self.clauses)):
                        q_ = subst(theta, q)
                        if all([unify(x, q_, {}) is None for x in self.clauses + new]):
                            new.append(q_)
                            phi = unify(q_, alpha, {})

                            if phi is not None:
                                yield phi
            if not new:
                break
            for clause in new:
                KB.tell(clause)
        return None

# kb = kb_engine()
# kb.tell('GPU(Cloud)')
# a = list(kb.fol_fc_ask('GPU(x)'))
# print(a[0], type(a[0]))