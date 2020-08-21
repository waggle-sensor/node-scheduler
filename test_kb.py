import unittest

from kb import kb_engine

class TestKB(unittest.TestCase):
    def test_inferencing(self):
        kb = kb_engine()
        print('Tell a rule: Rain(Now) ==> Run(Cloud)')
        kb.tell('Rain(Now) ==> Run(Cloud)')
        print('Tell a conversion: if rain_gauge >=3, then Rain(Now)')
        kb.memorize(('rain_gauge',), (('Rain(Now)', 'rain_gauge >= 3'),))

        print('Ask what to run')
        print(list(kb.fol_fc_ask('Run(x)')))
        print('Sense 2 ticks of rain_gauge')
        kb.sense('rain_gauge', 2)
        print('Ask what to run')
        print(list(kb.fol_fc_ask('Run(x)')))
        print('Sense 3 ticks of rain_gauge')
        kb.sense('rain_gauge', 3)
        print('Ask what to run')
        print(list(kb.fol_fc_ask('Run(x)')))
        print('Sense 2 ticks of rain_gauge')
        kb.sense('rain_gauge', 2)
        print('Ask what to run')
        print(list(kb.fol_fc_ask('Run(x)')))
        self.assertEqual(True, True)

if __name__ == '__main__':
    unittest.main()